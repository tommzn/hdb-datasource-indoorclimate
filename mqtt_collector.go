package indoorclimate

import (
	"context"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// NewMqttCollector returns a new collector listen for MQTT messages to obtain indoor climate data.
func NewMqttCollector(conf config.Config, logger log.Logger) *MqttCollector {

	mqttCollector := &MqttCollector{
		logger:        logger,
		publisher:     []Publisher{},
		measurements:  make(chan IndoorClimateMeasurement, 10),
		subscriptions: []MqttSubscriptionConfig{},
	}

	mqttCollector.mqttOptions = mqttCollector.mqttClientOptionsFromConfig(conf)
	mqttCollector.assignSubscriptions(conf)
	return mqttCollector
}

// Run start indoor climate collector. It's an infinite loop until given context is canceled.
func (collector *MqttCollector) Run(ctx context.Context) error {

	client, connectError := collector.mqttConnect()
	if connectError != nil {
		collector.logger.Error("Unable to connect to MQTT broker, reason: ", connectError)
		return connectError
	}

	go collector.subscribe(client, ctx)

	<-ctx.Done()
	client.Disconnect(250)
	return nil
}

// Subscribe to MQTT topics and handle incoming message by device plugins to extract indoor climate data.
func (collector *MqttCollector) subscribe(client mqtt.Client, ctx context.Context) {

	for _, subscription := range collector.subscriptions {
		subscription.Plugin.SetMeasurementChannel(collector.measurements)
		token := client.Subscribe(subscription.Topic, 1, subscription.Plugin.MessageHandler)
		token.Wait()
		if token.Error() != nil {
			collector.logger.Errorf("Unable to subscribe to topic; %s, reason: ", subscription.Topic, token.Error())
		}
	}

	for {
		select {
		case measurement := <-collector.measurements:
			collector.logger.Debugf("Measurement obtained: %+v", measurement)
			for _, publisher := range collector.publisher {
				if err := publisher.SendMeasurement(measurement); err != nil {
					collector.logger.Error(err)
				}
			}
		case <-ctx.Done():
			collector.logger.Debug("Stop sensor data collection: ", ctx.Err())
			return
		}
	}
}

// MqttConnect creates a MQTT client and try to connect to given broker.
func (collector *MqttCollector) mqttConnect() (mqtt.Client, error) {

	client := mqtt.NewClient(collector.mqttOptions)
	token := client.Connect()
	token.Wait()
	return client, token.Error()
}

// AooendTarget will append passed target to internal publisher list.
func (collector *MqttCollector) AppendTarget(newTarget Publisher) {
	collector.publisher = append(collector.publisher, newTarget)
}

// ConnectHandler logs status message after connection to MQTT broker has been established.
func (collector *MqttCollector) connectHandler(client mqtt.Client) {
	collector.logger.Info("Connected")
}

// ConnectLostHandler logs errors in case connection to MQTT broker get lost.
func (collector *MqttCollector) connectLostHandler(client mqtt.Client, err error) {
	collector.logger.Errorf("Connect lost: %v", err)
}

// ExtractFromConfig is a helper to get config data from given map.
func extractFromConfig(conf map[string]string, key string) *string {
	if val, ok := conf[key]; ok {
		return config.AsStringPtr(val)
	} else {
		return nil
	}
}

// NewDevicePlugin create a new device plugin for given key. If an unknow kex is passed nil is returned.
func newDevicePlugin(pluginKey *string) DevicePlugin {

	if pluginKey == nil {
		return nil
	}

	switch DevicePluginKey(*pluginKey) {
	default:
		return nil
	}
}

// MqttClientOptionsFromConfig extracts MQTT client settings from given config.
func (collector *MqttCollector) mqttClientOptionsFromConfig(conf config.Config) *mqtt.ClientOptions {

	broker := conf.Get("mqtt.broker", config.AsStringPtr("localhost"))
	port := conf.GetAsInt("mqtt.port", config.AsIntPtr(1883))
	options := mqtt.NewClientOptions()
	options.SetOrderMatters(false)
	options.AddBroker(fmt.Sprintf("tcp://%s:%d", *broker, *port))
	options.SetClientID("indoorclimate_collector")
	options.OnConnect = collector.connectHandler
	options.OnConnectionLost = collector.connectLostHandler
	return options
}

// AssignSubscriptions parse given config to extract topics this client should subscribe and created a corresponding device plugin.
func (collector *MqttCollector) assignSubscriptions(conf config.Config) {

	subscriptions := conf.GetAsSliceOfMaps("mqtt.subscriptions")
	for _, subscription := range subscriptions {

		topic := extractFromConfig(subscription, "topic")
		pluginKey := extractFromConfig(subscription, "plugin")
		devicePlugin := newDevicePlugin(pluginKey)
		if topic != nil && devicePlugin != nil {
			collector.subscriptions = append(collector.subscriptions, MqttSubscriptionConfig{Topic: *topic, Plugin: devicePlugin})
		}
	}
}
