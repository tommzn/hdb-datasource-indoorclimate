package indoorclimate

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

func NewMqttLivenessObserver(conf config.Config, logger log.Logger) *MqttLivenessObserver {

	livenessTopic := conf.Get("mqtt.liveness.topic", config.AsStringPtr("/hdb/liveness"))
	schedule := conf.GetAsDuration("mqtt.liveness.schedule", config.AsDurationPtr(1*time.Minute))
	wait := conf.GetAsDuration("mqtt.liveness.wait", config.AsDurationPtr(5*time.Second))
	observer := &MqttLivenessObserver{
		logger:        logger,
		livenessTopic: *livenessTopic,
		schedule:      *schedule,
		wait:          *wait,
		probeChan:     make(chan string, 1),
	}
	observer.mqttOptions = observer.mqttClientOptionsFromConfig(conf)
	return observer
}

func (observer *MqttLivenessObserver) Run(ctx context.Context) error {

	defer observer.logger.Flush()

	client, connectError := observer.mqttConnect()
	if connectError != nil {
		observer.logger.Error("Unable to connect to MQTT broker, reason: ", connectError)
		return connectError
	}

	token := client.Subscribe(observer.livenessTopic, 1, observer.MessageHandler)
	token.Wait()
	if subscribeErr := token.Error(); subscribeErr != nil {
		observer.logger.Error("Unable to subscribe to liveness topic, reason: ", subscribeErr)
		return subscribeErr
	}

	errorChan := make(chan error, 1)
	ticker := time.NewTicker(observer.schedule)
	observer.logger.Debugf("Linveness schedule: %s", observer.schedule)
	observer.logger.Debugf("Linveness wait: %s", observer.wait)
	go func() {

		for {
			select {
			case <-ticker.C:
				err := observer.liveness(client)
				if err != nil {
					errorChan <- err
					return
				}

			case <-ctx.Done():
				errorChan <- nil
			}
		}

	}()

	err := <-errorChan
	client.Disconnect(250)
	return err
}

func (observer *MqttLivenessObserver) liveness(client mqtt.Client) error {

	defer observer.logger.Flush()

	message := randomMessage(32)
	token := client.Publish(observer.livenessTopic, 1, false, message)
	token.Wait()
	if publishErr := token.Error(); publishErr != nil {
		observer.logger.Error("Unable to public liveness probe, reason: ", publishErr)
		return nil
	}
	observer.logger.Status("Liveness message send.")

	timer := time.NewTimer(observer.wait)
	select {
	case receivedMsg := <-observer.probeChan:
		if receivedMsg == message {
			observer.logger.Status("Liveness probe passed.")
		} else {
			observer.logger.Error("Receive invalid liveness probe.")
		}
	case <-timer.C:
		observer.logger.Error("Liveness timeout.")
	}
	return nil
}

func (observer *MqttLivenessObserver) MessageHandler(mclient mqtt.Client, msg mqtt.Message) {
	observer.logger.Statusf("Liveness message received. Topic: %s, MEssage: '%s'", msg.Topic(), msg.Payload())
	observer.probeChan <- string(msg.Payload())
}

// MqttConnect creates a MQTT client and try to connect to given broker.
func (observer *MqttLivenessObserver) mqttConnect() (mqtt.Client, error) {

	client := mqtt.NewClient(observer.mqttOptions)
	token := client.Connect()
	token.Wait()
	return client, token.Error()
}

// ConnectHandler logs status message after connection to MQTT broker has been established.
func (observer *MqttLivenessObserver) connectHandler(client mqtt.Client) {
	observer.logger.Info("Connected")
}

// ConnectLostHandler logs errors in case connection to MQTT broker get lost.
func (observer *MqttLivenessObserver) connectLostHandler(client mqtt.Client, err error) {
	observer.logger.Errorf("Connect lost: %v", err)
}

// MqttClientOptionsFromConfig extracts MQTT client settings from given config.
func (observer *MqttLivenessObserver) mqttClientOptionsFromConfig(conf config.Config) *mqtt.ClientOptions {

	broker := conf.Get("mqtt.broker", config.AsStringPtr("localhost"))
	port := conf.GetAsInt("mqtt.port", config.AsIntPtr(1883))
	options := mqtt.NewClientOptions()
	options.SetOrderMatters(false)
	options.AddBroker(fmt.Sprintf("tcp://%s:%d", *broker, *port))
	options.SetClientID("indoorclimate_collector_liveness")
	options.OnConnect = observer.connectHandler
	options.OnConnectionLost = observer.connectLostHandler
	return options
}

func randomMessage(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}
