package indoorclimate

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Publisher sends given measuremnts to different targets.
type Publisher interface {

	// SendMeasurement will start to transfer passed measurement to a target.
	SendMeasurement(IndoorClimateMeasurement) error
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
