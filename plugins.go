package indoorclimate

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

type testPlugin struct {
	measurementChan chan<- IndoorClimateMeasurement
	logger          log.Logger
	deviceId        string
}

func newtestPlugin(logger log.Logger, deviceId *string) *testPlugin {
	if deviceId == nil {
		deviceId = config.AsStringPtr("Device01")
	}
	return &testPlugin{logger: logger, deviceId: *deviceId}
}

func (plugin *testPlugin) MessageHandler(client mqtt.Client, message mqtt.Message) {

	plugin.logger.Debugf("Received, Topic: %s, Message: %s", message.Topic(), message.Payload())

	var measurement IndoorClimateMeasurement
	switch string(message.Payload()) {
	case "temperature":
		measurement = IndoorClimateMeasurement{
			DeviceId:  plugin.deviceId,
			Timestamp: time.Now(),
			Type:      MEASUREMENTTYPE_TEMPERATURE,
			Value:     "24.7",
		}
	case "humidity":
		measurement = IndoorClimateMeasurement{
			DeviceId:  plugin.deviceId,
			Timestamp: time.Now(),
			Type:      MEASUREMENTTYPE_HUMIDITY,
			Value:     "53.9",
		}
	case "battery":
		measurement = IndoorClimateMeasurement{
			DeviceId:  plugin.deviceId,
			Timestamp: time.Now(),
			Type:      MEASUREMENTTYPE_BATTERY,
			Value:     "87",
		}
	default:
		measurement = IndoorClimateMeasurement{
			DeviceId:  plugin.deviceId,
			Timestamp: time.Now(),
			Type:      MEASUREMENTTYPE_TEMPERATURE,
			Value:     "12.3",
		}
	}

	plugin.logger.Debugf("Generated Measurement: %+v", measurement)
	plugin.measurementChan <- measurement
}

// SetMeasurementChannel assign a channel measuremnts should be written to.
func (plugin *testPlugin) SetMeasurementChannel(channel chan<- IndoorClimateMeasurement) {
	plugin.measurementChan = channel
}
