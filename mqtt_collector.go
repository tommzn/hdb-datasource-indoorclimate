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

// AppendSubscription will append passed subscription to internal subscription list.
func (collector *MqttCollector) AppendSubscription(subscription MqttSubscriptionConfig) {
	collector.subscriptions = append(collector.subscriptions, subscription)
}

// ConnectHandler logs status message after connection to MQTT broker has been established.
func (collector *MqttCollector) connectHandler(client mqtt.Client) {
	collector.logger.Info("Connected")
}

// ConnectLostHandler logs errors in case connection to MQTT broker get lost.
func (collector *MqttCollector) connectLostHandler(client mqtt.Client, err error) {
	collector.logger.Errorf("Connect lost: %v", err)
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
