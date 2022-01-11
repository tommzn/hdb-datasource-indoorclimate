package indoorclimate

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	metrics "github.com/tommzn/go-metrics"
	secrets "github.com/tommzn/go-secrets"
	events "github.com/tommzn/hdb-events-go"
)

// MQTT_CLIENT_ID to be used at MQTT connections
const MQTT_CLIENT_ID = "indoorclimate_consumer"

// New returns a new MQTT client to subscribe to topics and process messages.
func New(conf config.Config, logger log.Logger, secretsManager secrets.SecretsManager) Collector {
	return &MqttClient{
		conf:            conf,
		logger:          logger,
		secretsManager:  secretsManager,
		targets:         []MessageTarget{newLogTarget(logger)},
		metricPublisher: metrics.NewTimestreamPublisher(conf, logger),
	}
}

// AppendMessageTarget add passed target to the internal stack.
func (client *MqttClient) AppendMessageTarget(target MessageTarget) {
	client.targets = append(client.targets, target)
}

// run creates a MQTT client, connects to a given brokcer and listen for indoor climate data.
// Will run until passed context has been canceled.
func (client *MqttClient) Run(ctx context.Context) error {

	defer client.logger.Flush()
	defer client.metricPublisher.Flush()

	filters := client.mqttTopicFilters()
	opts := client.mqttOptions()
	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		client.logger.Errorf("Unable to connect to broker (%s), reason: %s", brokerList(opts.Servers), token.Error())
		return token.Error()
	}

	if token := mqttClient.SubscribeMultiple(filters, client.processMessage); token.Wait() && token.Error() != nil {
		client.logger.Errorf("Unable to subsribe to topics, reason: %s", token.Error())
		return token.Error()
	}

	<-ctx.Done()
	client.logger.Debug("Stop message consuming!")
	if mqttClient.IsConnected() {
		mqttClient.Disconnect(0)
	}
	return nil
}

// connectHandler is called after successful connection to a MQTT broker.
func (client *MqttClient) connectHandler(mqttClient mqtt.Client) {
	client.logger.Info("Connected to MQTT broker.")
	client.logger.Flush()
}

// connectionLostHandler is called if connection to a MQTT broker get lost.
func (client *MqttClient) connectionLostHandler(mqttClient mqtt.Client, err error) {
	opts := mqttClient.OptionsReader()
	client.logger.Infof("Connection to MQTT broker lost: %s, reason: %s", brokerList(opts.Servers()), err.Error())
	client.logger.Flush()
}

// mqttTopicFilters adds a prefix to consumed topics if defined.
func (client *MqttClient) mqttTopicFilters() map[string]byte {

	filters := make(map[string]byte)
	topics := topicsToSubsrcibe(client.conf.Get("mqtt.topic_prefix", nil))
	for _, topics := range topics {
		filters[topics] = 0
	}
	return filters
}

// processMessage is called after a new message has been received from MQTT topic.
// It will convert reeived data to indoor climate data and calls all message targets in local stack in sequence.
func (client *MqttClient) processMessage(mqttClient mqtt.Client, message mqtt.Message) {

	defer client.logger.Flush()
	client.logger.Debugf("Receive: Topic: %s, Payload: %s", message.Topic(), message.Payload())

	deviceId := extractDeviceId(message.Topic())
	if deviceId == nil {
		client.logger.Error("Unable to get device id from topic: ", message.Topic())
		return
	}

	measurementTypeString := extractMeasurementType(message.Topic())
	if measurementTypeString == nil {
		client.logger.Error("Unable to extract measurement type from topic: ", message.Topic())
		return
	}

	measurementType, ok := events.MeasurementType_value[strings.ToUpper(*measurementTypeString)]
	if !ok {
		client.logger.Error("Unable to get measurement type for: ", measurementTypeString)
		return
	}

	indoorClimate := events.IndoorClimate{
		DeviceId:  *deviceId,
		Timestamp: timestamppb.New(time.Now()),
		Type:      events.MeasurementType(measurementType),
		Value:     string(message.Payload()),
	}
	client.metricPublisher.Send(createMeasurement(indoorClimate))
	for _, target := range client.targets {
		target.Send(indoorClimate)
	}
}

// mqttOptions defines options to connect to a MQTT broker.
func (client *MqttClient) mqttOptions() *mqtt.ClientOptions {

	broker := client.conf.Get("mqtt.broker", config.AsStringPtr("localhost"))
	port := client.conf.GetAsInt("mqtt.port", config.AsIntPtr(1883))
	opts := mqtt.NewClientOptions()
	brokerUrl := fmt.Sprintf("tcp://%s:%d", *broker, *port)
	opts.AddBroker(brokerUrl)
	clientId := MQTT_CLIENT_ID + "_" + randStringBytes(5)
	opts.SetClientID(clientId)
	opts.AutoReconnect = true
	client.logger.Debugf("MQTT Opts: %s", opts)
	opts.CredentialsProvider = client.credentialsProvider
	opts.OnConnect = client.connectHandler
	opts.OnConnectionLost = client.connectionLostHandler
	return opts
}

// credentialsProvider will return user name and password if provided by local secrets mananger.
func (client *MqttClient) credentialsProvider() (username string, password string) {

	mqttUser, _ := client.secretsManager.Obtain("TSL_MQTT_USER")
	mqttPassword, _ := client.secretsManager.Obtain("TSL_MQTT_PASSWORD")
	if mqttUser != nil && mqttPassword != nil {
		username = *mqttUser
		password = *mqttPassword
	}
	return username, password
}

// brokerList converts a list of broker urls to a single string.
func brokerList(urls []*url.URL) string {
	broker := []string{}
	for _, url := range urls {
		broker = append(broker, url.Host)
	}
	return strings.Join(broker, ",")
}

func createMeasurement(indoorClimate events.IndoorClimate) metrics.Measurement {
	return metrics.Measurement{
		MetricName: "hdb-datasource-indoorclimate",
		Tags: []metrics.MeasurementTag{
			metrics.MeasurementTag{
				Name:  "deviceid",
				Value: indoorClimate.DeviceId,
			},
			metrics.MeasurementTag{
				Name:  "type",
				Value: indoorClimate.Type.String(),
			},
		},
		Values: []metrics.MeasurementValue{
			metrics.MeasurementValue{
				Name:  "count",
				Value: 1,
			},
			metrics.MeasurementValue{
				Name:  indoorClimate.Type.String(),
				Value: indoorClimate.Value,
			},
		},
	}
}
