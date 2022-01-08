package main

import (
	"context"
	"os"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"

	core "github.com/tommzn/hdb-datasource-core"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

func main() {

	ctx := context.Background()
	collector, err := bootstrap(ctx)
	if err != nil {
		panic(err)
	}
	collector.Run(ctx)
}

// bootstrap loads config and creates a new scheduled collector with a exchangerate datasource.
func bootstrap(ctx context.Context) (core.Collector, error) {

	secretsManager := newSecretsManager()
	conf := loadConfig()
	logger := newLogger(conf, secretsManager, ctx)
	datasource := indoorclimate.New(conf, logger, secretsManager)
	datasource.AppendMessageTarget(indoorclimate.NewSqsTarget(conf, logger))
	return core.NewContinuousCollector(datasource, logger), nil
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
	logContextValues[log.LogCtxNamespace] = "hdb-datasource-indoorclimate"
	if node, ok := os.LookupEnv("K8S_NODE_NAME"); ok {
		logContextValues[log.LogCtxK8sNode] = node
	}
	if pod, ok := os.LookupEnv("K8S_POD_NAME"); ok {
		logContextValues[log.LogCtxK8sPod] = pod
	}
	logger.WithContext(log.LogContextWithValues(ctx, logContextValues))
	return logger
}
