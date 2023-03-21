package plugins

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/tommzn/go-log"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

type ShellyHTPlugin struct {
	measurementChan chan<- indoorclimate.IndoorClimateMeasurement
	logger          log.Logger
}

type ShellyHT_Event struct {
	Method string `json:"method"`
}

type ShellyHT_NotifyFullStatus struct {
	Device string          `json:"src"`
	Params ShellyHT_Params `json:"params"`
}

type ShellyHT_Params struct {
	DevicePower0 ShellyHT_DevicePower `json:"devicepower:0"`
	Temperature0 ShellyHT_Temperature `json:"temperature:0"`
	Humidity0    ShellyHT_Humidity    `json:"humidity:0"`
}

type ShellyHT_Humidity struct {
	Value float64 `json:"rh"`
}

type ShellyHT_Temperature struct {
	ValueC float64 `json:"tC"`
	ValueF float64 `json:"tF"`
}

type ShellyHT_DevicePower struct {
	Battery  ShellyHT_BatteryPower  `json:"battery"`
	External ShellyHT_ExternalPower `json:"external"`
}

type ShellyHT_BatteryPower struct {
	Value float64 `json:"percent"`
}

type ShellyHT_ExternalPower struct {
	Present bool `json:"present"`
}

// NewMockPlugin creates a new mock plugin.
func NewShellyHTPlugin(logger log.Logger) *ShellyHTPlugin {
	return &ShellyHTPlugin{logger: logger}
}

type SellyHT_EventType string

const (
	SHELLYHT_NOTIFYFULLSTATUS SellyHT_EventType = "NotifyFullStatus"
)

// MessageHandler process message received from MQTT topic and pulishes static indoor climate data depending on message payload.
func (plugin *ShellyHTPlugin) MessageHandler(client mqtt.Client, message mqtt.Message) {

	plugin.logger.Debugf("Received, Topic: %s, Message: %s", message.Topic(), message.Payload())

	eventType := getEventType(message.Payload())
	plugin.logger.Statusf("Event Type: %s", eventType)
	if eventType == SHELLYHT_NOTIFYFULLSTATUS {

		var event ShellyHT_NotifyFullStatus
		if err := json.Unmarshal(message.Payload(), &event); err == nil {

			plugin.publishTemperatureMeasurement(event.Device, event.Params.Temperature0)
			plugin.publishHumidityMeasurement(event.Device, event.Params.Humidity0)
			plugin.publishBatteryMeasurement(event.Device, event.Params.DevicePower0)
		}
	}
}

// SetMeasurementChannel assign a channel measuremnts should be written to.
func (plugin *ShellyHTPlugin) SetMeasurementChannel(channel chan<- indoorclimate.IndoorClimateMeasurement) {
	plugin.measurementChan = channel
}

// PublishTemperatureMeasurement creates and publishes a temperature measurement.
func (plugin *ShellyHTPlugin) publishTemperatureMeasurement(device string, temp ShellyHT_Temperature) {
	measurement := indoorclimate.IndoorClimateMeasurement{
		DeviceId:  device,
		Timestamp: time.Now(),
		Type:      indoorclimate.MEASUREMENTTYPE_TEMPERATURE,
		Value:     fmt.Sprintf("%.1f", temp.ValueC),
	}
	plugin.logger.Debugf("Generated Temparature Measurement: %+v", measurement)
	plugin.measurementChan <- measurement
}

// PublishHumidityMeasurement creates and publishes a humidity measurement.
func (plugin *ShellyHTPlugin) publishHumidityMeasurement(device string, humidity ShellyHT_Humidity) {
	measurement := indoorclimate.IndoorClimateMeasurement{
		DeviceId:  device,
		Timestamp: time.Now(),
		Type:      indoorclimate.MEASUREMENTTYPE_HUMIDITY,
		Value:     fmt.Sprintf("%.1f", humidity.Value),
	}
	plugin.logger.Debugf("Generated Humidity Measurement: %+v", measurement)
	plugin.measurementChan <- measurement
}

// PublishBatteryMeasurement creates and publishes a battery measurement.
func (plugin *ShellyHTPlugin) publishBatteryMeasurement(device string, power ShellyHT_DevicePower) {
	measurement := indoorclimate.IndoorClimateMeasurement{
		DeviceId:  device,
		Timestamp: time.Now(),
		Type:      indoorclimate.MEASUREMENTTYPE_BATTERY,
		Value:     fmt.Sprintf("%.1f", power.Battery.Value),
	}
	plugin.logger.Debugf("Generated Bettery Measurement: %+v", measurement)
	plugin.measurementChan <- measurement
}

// GetEventType extracts Shelly event type from given MQTT message.
func getEventType(messagePayload []byte) SellyHT_EventType {
	var event ShellyHT_Event
	json.Unmarshal(messagePayload, &event)
	return SellyHT_EventType(event.Method)
}
