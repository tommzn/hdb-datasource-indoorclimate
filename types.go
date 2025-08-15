package indoorclimate

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/tommzn/go-log"
)

// DevicePluginKey defines a plugin for devices of a specific manufacturer.
type DevicePluginKey string

const (
	PLUGIN_SHELLY         DevicePluginKey = "shelly"
	PLUGIN_LOGGER         DevicePluginKey = "logger"
	PLUGIN_HOME_ASSISTANT DevicePluginKey = "homeassistant"
)

// MeasurementType is a indoor climate date type, e.g. temperature.
type MeasurementType string

const (
	MEASUREMENTTYPE_TEMPERATURE MeasurementType = "temperature"
	MEASUREMENTTYPE_HUMIDITY    MeasurementType = "humidity"
	MEASUREMENTTYPE_BATTERY     MeasurementType = "battery"
)

// IndoorClimateMeasurement is a metric read from a sensor device.
type IndoorClimateMeasurement struct {
	DeviceId  string
	Timestamp time.Time
	Type      MeasurementType
	Value     string
}

// Characteristic is a songle sensor value.
type Characteristic struct {
	uuid            string
	measurementType MeasurementType
}

// MqttCollector subcribes to MQTT topics to extract indoor climate date from published messages.
type MqttCollector struct {
	logger        log.Logger
	publisher     []Publisher
	measurements  chan IndoorClimateMeasurement
	subscriptions []MqttSubscriptionConfig
	mqttOptions   *mqtt.ClientOptions
}

// MqttSubscriptionConfig define a MQTT topic and it's message handler plugin.
type MqttSubscriptionConfig struct {
	Topic  string
	Plugin DevicePlugin
}

// MqttLivenessObserver can be use to observer a MQTT server.
type MqttLivenessObserver struct {
	logger        log.Logger
	mqttOptions   *mqtt.ClientOptions
	livenessTopic string
	schedule      time.Duration
	wait          time.Duration
	probeChan     chan string
}
