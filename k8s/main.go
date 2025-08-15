package main

import (
	"context"
	"os"
	"sync"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"

	hdbcore "github.com/tommzn/hdb-core"
	core "github.com/tommzn/hdb-datasource-core"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
	plugins "github.com/tommzn/hdb-datasource-indoorclimate/plugins"
	targets "github.com/tommzn/hdb-datasource-indoorclimate/targets"
)

func main() {

	ctx := context.Background()
	collector, LivenessObserver, err := bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	if LivenessObserver != nil {
		go LivenessObserver.Run(ctx, wg)
	}
	collector.Run(ctx)
	wg.Wait()
}

// bootstrap loads config and creates a new scheduled collector with a exchangerate datasource.
func bootstrap(ctx context.Context) (core.Collector, hdbcore.Runable, error) {

	secretsManager := newSecretsManager()
	conf := loadConfig()
	logger := newLogger(conf, secretsManager, ctx)
	datasource := indoorclimate.NewMqttCollector(conf, logger, newSecretsManager())
	if queue := conf.Get("hdb.queue", nil); queue != nil {
		datasource.AppendTarget(targets.NewSqsTarget(conf, logger))
	}
	if timestreamTable := conf.Get("aws.timestream.table", nil); timestreamTable != nil {
		datasource.AppendTarget(targets.NewTimestreamTarget(conf, logger))
	}
	if topicArn := conf.Get("hdb.topic.arn", nil); topicArn != nil {
		datasource.AppendTarget(targets.NewSnsTarget(conf))
	}
	if bucket := conf.Get("aws.s3.bucket", nil); bucket != nil {
		if s3Target, err := targets.NewS3Target(conf); err == nil {
			datasource.AppendTarget(s3Target)
		}
	}
	subsriptions := SubsriptionsFromConfig(conf, logger)
	for _, subsription := range subsriptions {
		datasource.AppendSubscription(subsription)
	}

	var livenessObserver hdbcore.Runable
	livenessTopic := conf.Get("mqtt.lineness.topic", nil)
	if livenessTopic == nil {
		livenessObserver = indoorclimate.NewMqttLivenessObserver(conf, logger)
	}

	return core.NewContinuousCollector(datasource, logger), livenessObserver, nil
}

// loadConfig from config file.
func loadConfig() config.Config {

	configSource, err := config.NewS3ConfigSourceFromEnv()
	if err != nil {
		panic(err)
	}

	conf, err := configSource.Load()
	if err != nil {
		panic(err)
	}
	return conf
}

// newSecretsManager retruns a new container secrets manager
func newSecretsManager() secrets.SecretsManager {
	secretsManager := secrets.NewDockerecretsManager("/run/secrets/token")
	secrets.ExportToEnvironment([]string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"}, secretsManager)
	return secretsManager
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager, ctx context.Context) log.Logger {
	logger := log.NewLoggerFromConfig(conf, secretsMenager)
	logContextValues := make(map[string]string)
	logContextValues[log.LogCtxNamespace] = "hdb-datasource-indoorclimate-mqtt"
	if node, ok := os.LookupEnv("K8S_NODE_NAME"); ok {
		logContextValues[log.LogCtxK8sNode] = node
	}
	if pod, ok := os.LookupEnv("K8S_POD_NAME"); ok {
		logContextValues[log.LogCtxK8sPod] = pod
	}
	logger.WithContext(log.LogContextWithValues(ctx, logContextValues))
	return logger
}

// SubsriptionsFromConfig parse given config to extract topics this client should subscribe and created a corresponding device plugin.
func SubsriptionsFromConfig(conf config.Config, logger log.Logger) []indoorclimate.MqttSubscriptionConfig {

	listOfSubsriptions := []indoorclimate.MqttSubscriptionConfig{}
	subscriptions := conf.GetAsSliceOfMaps("mqtt.subscriptions")
	for _, subscription := range subscriptions {

		topic := extractFromConfig(subscription, "topic")
		pluginKey := extractFromConfig(subscription, "plugin")
		devicePlugin := newDevicePlugin(pluginKey, logger)
		if topic != nil && devicePlugin != nil {
			listOfSubsriptions = append(listOfSubsriptions, indoorclimate.MqttSubscriptionConfig{Topic: *topic, Plugin: devicePlugin})
		}
	}
	return listOfSubsriptions
}

// NewDevicePlugin create a new device plugin for given key. If an unknow kex is passed nil is returned.
func newDevicePlugin(pluginKey *string, logger log.Logger) indoorclimate.DevicePlugin {

	if pluginKey == nil {
		return nil
	}

	switch indoorclimate.DevicePluginKey(*pluginKey) {
	case indoorclimate.PLUGIN_SHELLY:
		return plugins.NewShellyHTPlugin(logger)
	case indoorclimate.PLUGIN_LOGGER:
		return plugins.NewLoggerPlugin(logger)
	case indoorclimate.PLUGIN_HOME_ASSISTANT:
		return plugins.NewHomeAssistantPlugin(logger)
	default:
		return nil
	}
}

// ExtractFromConfig is a helper to get config data from given map.
func extractFromConfig(conf map[string]string, key string) *string {
	if val, ok := conf[key]; ok {
		return config.AsStringPtr(val)
	} else {
		return nil
	}
}
