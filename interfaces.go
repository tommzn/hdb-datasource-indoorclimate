package indoorclimate

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Publisher sends given measuremnts to different targets.
type Publisher interface {

	// SendMeasurement will start to transfer passed measurement to a target.
	SendMeasurement(IndoorClimateMeasurement) error
}

// SensorDevice represents a device to fetch indoor cliamte data.
type SensorDevice interface {

	// Returns the id of current sensor device.
	Id() string

	// Connect will try to connect to a device and will return with an error if failing.
	Connect() error

	// Disconnect will try to disconnect from current device and returns with an error if it fails.
	Disconnect() error

	// ReadValue will try to read measurment value for given characteristics.
	ReadValue(string) ([]byte, error)
}

// MqttSubscription is used to handle messages received from subscribes MQTT topics.
type MqttSubscription interface {

	// MessageHandler process received data received from a mqtt topic
	MessageHandler(mqtt.Client, mqtt.Message)
}

// DevicePlugin is used to subscribe to MQTT topics and extract measurements from published data.
type DevicePlugin interface {
	MqttSubscription

	// SetMeasurementChannel assigns a channel plugin should write extracted indoor climate measurements to.
	SetMeasurementChannel(chan<- IndoorClimateMeasurement)
}
