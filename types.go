package indoorclimate

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	btdevice "github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/tommzn/go-log"
	utils "github.com/tommzn/go-utils"
)

// DevicePluginKey defines a plugin for devices of a specific manufacturer.
type DevicePluginKey string

const (
	PLUGIN_SHELLY DevicePluginKey = "shelly"
)

// MeasurementType is a indoor climate date type, e.g. temperature.
type MeasurementType string

const (
	MEASUREMENTTYPE_TEMPERATURE MeasurementType = "temperature"
	MEASUREMENTTYPE_HUMIDITY    MeasurementType = "humidity"
	MEASUREMENTTYPE_BATTERY     MeasurementType = "battery"
)

// IndoorClimateSensor is used to fetch tem eprature, humidiy and bettery status
// from a Xiaomi Mijia (LYWSD03MMC) indoor climate sensor.
type IndoorClimateSensor struct {
	adapterId string
	deviceId  string
	device    *btdevice.Device1
}

// IndoorClimateMeasurement is a metric read from a sensor device.
type IndoorClimateMeasurement struct {
	DeviceId  string
	Timestamp time.Time
	Type      MeasurementType
	Value     string
}

// LogPublisher will log indoor climate measuremnts.
type LogPublisher struct {
	logger log.Logger
}

// Characteristic is a songle sensor value.
type Characteristic struct {
	uuid            string
	measurementType MeasurementType
}

// SensorDataCollector will try to fetch temperature, humidity and bettery status
// from a given list of sensors.
type SensorDataCollector struct {
	logger          log.Logger
	devices         []SensorDevice
	characteristics []Characteristic
	publisher       []Publisher
	retryCount      int
	schedule        *time.Duration
	errorStack      *utils.ErrorStack
	done            chan struct{}
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
