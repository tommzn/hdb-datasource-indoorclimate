package plugins

import (
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/tommzn/go-log"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

// MeasurementType defines the type of measurement being reported by a sensor.
type MeasurementType string

const (
	// MeasurementTypeTemperature indicates a temperature measurement.
	MeasurementTypeTemperature MeasurementType = "temperature"
	// MeasurementTypeHumidity indicates a humidity measurement.
	MeasurementTypeHumidity = "humidity"
	// MeasurementTypeBattery indicates a battery level measurement.
	MeasurementTypeBattery = "battery"
)

// HomeAssistantSensor represents the structure of a Home Assistant sensor data point.
type HomeAssistantSensor struct {
	// EntityID is the unique identifier for the sensor entity in Home Assistant.
	EntityID string `json:"entity_id"`
	// Name is the user-friendly name of the sensor.
	Name string `json:"name"`
	// Room indicates the location of the sensor within the home.
	Room string `json:"room"`
	// Time is the timestamp when the measurement was taken.
	Time time.Time `json:"time"`
	// TypeOfMeasurement specifies the kind of data being reported (e.g., temperature, humidity).
	TypeOfMeasurement MeasurementType `json:"type_of_measurement"`
	// Unit is the unit of measurement for the value (e.g., °C, %).
	Unit string `json:"unit"`
	// Value is the actual measured value.
	Value string `json:"value"`
}

// HomeAssistantPlugin is a plugin for collecting indoor climate measurements from Home Assistant sensors.
type HomeAssistantPlugin struct {
	// measurementChan is the channel to send the collected indoor climate measurements to.
	measurementChan chan<- indoorclimate.IndoorClimateMeasurement
	// logger is the logger instance used for logging messages within the plugin.
	logger log.Logger
}

// NewHomeAssistantPlugin creates a new plugin.
func NewHomeAssistantPlugin(logger log.Logger) *HomeAssistantPlugin {
	return &HomeAssistantPlugin{logger: logger}
}

// MessageHandler processes incoming MQTT messages, attempting to parse them as Home Assistant sensor data.
// It then converts the parsed data into an IndoorClimateMeasurement and sends it to the measurement channel.
func (plugin *HomeAssistantPlugin) MessageHandler(client mqtt.Client, message mqtt.Message) {
	plugin.logger.Debugf("[HomeAssistantPlugin] Measurement received, Topic: %s, Message: %s", message.Topic(), message.Payload())

	var event HomeAssistantSensor
	json.Unmarshal(message.Payload(), &event)
	if event.EntityID == "" {
		plugin.logger.Error("[HomeAssistantPlugin] Invalid payload received: ", message.Payload())
		return
	}

	var indoorClimateMeasurement indoorclimate.IndoorClimateMeasurement
	switch event.TypeOfMeasurement {
	case MeasurementTypeTemperature:
		indoorClimateMeasurement = asIndoorClimateMeasurement(event, indoorclimate.MEASUREMENTTYPE_TEMPERATURE)
	case MeasurementTypeHumidity:
		indoorClimateMeasurement = asIndoorClimateMeasurement(event, indoorclimate.MEASUREMENTTYPE_HUMIDITY)
	case MeasurementTypeBattery:
		indoorClimateMeasurement = asIndoorClimateMeasurement(event, indoorclimate.MEASUREMENTTYPE_BATTERY)

	default:
		plugin.logger.Errorf("[HomeAssistantPlugin] Unknown measurement type: %s", event.TypeOfMeasurement)
	}

	plugin.logger.Debugf("Generated Temparature Measurement: %+v", indoorClimateMeasurement)
	plugin.measurementChan <- indoorClimateMeasurement
}

// asIndoorClimateMeasurement converts a HomeAssistantSensor struct into an IndoorClimateMeasurement struct.
func asIndoorClimateMeasurement(event HomeAssistantSensor, measurementType indoorclimate.MeasurementType) indoorclimate.IndoorClimateMeasurement {
	return indoorclimate.IndoorClimateMeasurement{
		DeviceId:  event.EntityID,
		Timestamp: event.Time,
		Type:      measurementType,
		Value:     event.Value,
	}
}

// SetMeasurementChannel assign a channel measuremnts should be written to.
func (plugin *HomeAssistantPlugin) SetMeasurementChannel(channel chan<- indoorclimate.IndoorClimateMeasurement) {
	plugin.measurementChan = channel
}
