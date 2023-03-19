package plugins

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

// MockPlugin can be used for testing.
// it generates static indoor climate measuremnts based on received message.
//
//	Message		|	Indoor CLimate Data
//
// --------------------------------------
//
//	temperature	|	Temperature 24.7°C
//	humidity	|	Humidity 53.9%
//	battery		|	Capacity 87%
//	<default>	|	Temperature 12.3°C
type MockPlugin struct {
	measurementChan chan<- indoorclimate.IndoorClimateMeasurement
	logger          log.Logger
	deviceId        string
}

// NewMockPlugin creates a new mock plugin.
func NewMockPlugin(logger log.Logger, deviceId *string) *MockPlugin {
	if deviceId == nil {
		deviceId = config.AsStringPtr("Device01")
	}
	return &MockPlugin{logger: logger, deviceId: *deviceId}
}

// MessageHandler process message received from MQTT topic and pulishes static indoor climate data depending on message payload.
func (plugin *MockPlugin) MessageHandler(client mqtt.Client, message mqtt.Message) {

	plugin.logger.Debugf("Received, Topic: %s, Message: %s", message.Topic(), message.Payload())

	var measurement indoorclimate.IndoorClimateMeasurement
	switch string(message.Payload()) {
	case "temperature":
		measurement = indoorclimate.IndoorClimateMeasurement{
			DeviceId:  plugin.deviceId,
			Timestamp: time.Now(),
			Type:      indoorclimate.MEASUREMENTTYPE_TEMPERATURE,
			Value:     "24.7",
		}
	case "humidity":
		measurement = indoorclimate.IndoorClimateMeasurement{
			DeviceId:  plugin.deviceId,
			Timestamp: time.Now(),
			Type:      indoorclimate.MEASUREMENTTYPE_HUMIDITY,
			Value:     "53.9",
		}
	case "battery":
		measurement = indoorclimate.IndoorClimateMeasurement{
			DeviceId:  plugin.deviceId,
			Timestamp: time.Now(),
			Type:      indoorclimate.MEASUREMENTTYPE_BATTERY,
			Value:     "87",
		}
	default:
		measurement = indoorclimate.IndoorClimateMeasurement{
			DeviceId:  plugin.deviceId,
			Timestamp: time.Now(),
			Type:      indoorclimate.MEASUREMENTTYPE_TEMPERATURE,
			Value:     "12.3",
		}
	}

	plugin.logger.Debugf("Generated Measurement: %+v", measurement)
	plugin.measurementChan <- measurement
}

// SetMeasurementChannel assign a channel measuremnts should be written to.
func (plugin *MockPlugin) SetMeasurementChannel(channel chan<- indoorclimate.IndoorClimateMeasurement) {
	plugin.measurementChan = channel
}
