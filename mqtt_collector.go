package indoorclimate

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	oslog "log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	"github.com/tommzn/go-secrets"
)

const (
	mqtt_username string = "MQTT_USERNAME"
	mqtt_password string = "MQTT_PASSWORD"
)

// NewMqttCollector returns a new collector listen for MQTT messages to obtain indoor climate data.
func NewMqttCollector(conf config.Config, logger log.Logger, secretsManager secrets.SecretsManager) *MqttCollector {

	mqttCollector := &MqttCollector{
		logger:        logger,
		publisher:     []Publisher{},
		measurements:  make(chan IndoorClimateMeasurement, 10),
		subscriptions: []MqttSubscriptionConfig{},
	}

	mqttCollector.mqttOptions = mqttCollector.mqttClientOptionsFromConfig(conf, secretsManager)
	return mqttCollector
}

// Run start indoor climate collector. It's an infinite loop until given context is canceled.
func (collector *MqttCollector) Run(ctx context.Context) error {

	defer collector.logger.Flush()

	client, connectError := collector.mqttConnect()
	if connectError != nil {
		collector.logger.Error("Unable to connect to MQTT broker, reason: ", connectError)
		return connectError
	}
	collector.logger.Info("Connected to MQTT broker")
	oslog.Println("Connected to MQTT Broker.")

	go collector.subscribe(client, ctx)

	collector.logger.Info("Processing messages...")
	oslog.Println("Processing messages...")
	collector.logger.Flush()
	<-ctx.Done()
	collector.logger.Info("Process cancelation received, disconnect from MQTT.")
	oslog.Println("Process cancelation received, disconnect from MQTT.")
	client.Disconnect(250)
	return nil
}

// Subscribe to MQTT topics and handle incoming message by device plugins to extract indoor climate data.
func (collector *MqttCollector) subscribe(client mqtt.Client, ctx context.Context) {

	collector.logger.Debugf("Subscriptions; %d", len(collector.subscriptions))
	for _, subscription := range collector.subscriptions {
		subscription.Plugin.SetMeasurementChannel(collector.measurements)
		token := client.Subscribe(subscription.Topic, 1, subscription.Plugin.MessageHandler)
		token.Wait()
		if token.Error() != nil {
			collector.logger.Errorf("Unable to subscribe to topic; %s, reason: ", subscription.Topic, token.Error())
		} else {
			collector.logger.Debugf("Successful subscribed to topic %s", subscription.Topic)
			oslog.Println("Successful subscribed to topic ", subscription.Topic)

		}
	}
	collector.logger.Flush()

	for {
		select {
		case measurement := <-collector.measurements:
			collector.logger.Debugf("Measurement obtained: %+v", measurement)
			for _, publisher := range collector.publisher {
				if err := publisher.SendMeasurement(measurement); err != nil {
					collector.logger.Error(err)
				}
			}
			collector.logger.Flush()
		case <-ctx.Done():
			collector.logger.Debug("Stop sensor data collection: ", ctx.Err())
			oslog.Println("Canceled.")
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
func (collector *MqttCollector) mqttClientOptionsFromConfig(conf config.Config, secretsmanager secrets.SecretsManager) *mqtt.ClientOptions {

	broker := conf.Get("mqtt.broker", config.AsStringPtr("localhost"))
	port := conf.GetAsInt("mqtt.port", config.AsIntPtr(1883))
	protocol := conf.Get("mqtt.protocol", config.AsStringPtr("tcp"))
	options := mqtt.NewClientOptions()
	options.SetOrderMatters(false)
	options.AddBroker(fmt.Sprintf("%s://%s:%d", *protocol, *broker, *port))
	options.SetClientID("indoorclimate_collector")
	options.OnConnect = collector.connectHandler
	options.OnConnectionLost = collector.connectLostHandler

	if username, _ := secretsmanager.Obtain(mqtt_username); username != nil {
		options.SetUsername(*username)
	}
	if password, _ := secretsmanager.Obtain(mqtt_password); password != nil {
		options.SetPassword(*password)
	}

	if strings.HasPrefix(*protocol, "ssl") {
		certpool := x509.NewCertPool()
		tlsConfig := &tls.Config{
			RootCAs:            certpool,
			InsecureSkipVerify: false,
		}
		options.SetTLSConfig(tlsConfig)
	}
	collector.logger.Infof("MQTT Broker: %v", options.Servers)
	return options
}
