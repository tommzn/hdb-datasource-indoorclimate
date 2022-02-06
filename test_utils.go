package indoorclimate

import (
	"os"
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
	events "github.com/tommzn/hdb-events-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// loadConfigForTest loads test config.
func loadConfigForTest(fileName *string) config.Config {

	configFile := "fixtures/testconfig.yml"
	if fileName != nil {
		configFile = *fileName
	}
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

// secretsManagerForTest returns a default secrets manager for testing
func secretsManagerForTest() secrets.SecretsManager {
	return secrets.NewSecretsManager()
}

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

func messageConsumerForTest(configFile *string) *MqttClient {
	return New(loadConfigForTest(configFile), loggerForTest(), secretsManagerForTest()).(*MqttClient)
}

func messageForTest() mqtt.Message {
	return &mqttMessage{
		topic:   "iobroker/ble/0/a4:f3:e6:b8:d1:c6/temperature",
		payload: "23.4",
	}
}

func messageWithTopicForTest(topic string) mqtt.Message {
	return &mqttMessage{
		topic:   topic,
		payload: "23.4",
	}
}

func indoorCliamteDataForTest() events.IndoorClimate {
	return events.IndoorClimate{
		DeviceId:  "a1:a1:a1:a1:a1:a1",
		Timestamp: timestamppb.New(time.Now()),
		Type:      events.MeasurementType_TEMPERATURE,
		Value:     "23.5",
	}
}

// skipCI returns true if env variable CI is set
func skipCI(t *testing.T) {
	if _, isSet := os.LookupEnv("CI"); isSet {
		t.Skip("Skipping testing in CI environment")
	}
}

func waitGroupForTest(childCount int) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(childCount)
	return wg
}
