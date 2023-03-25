package plugins

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/tommzn/go-log"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

type LoggerPlugin struct {
	logger log.Logger
}

func NewLoggerPlugin(logger log.Logger) *LoggerPlugin {
	return &LoggerPlugin{logger: logger}
}

func (plugin *LoggerPlugin) MessageHandler(client mqtt.Client, message mqtt.Message) {
	plugin.logger.Debugf("Topic: %s, Message: %s", message.Topic(), message.Payload())
}

func (plugin *LoggerPlugin) SetMeasurementChannel(channel chan<- indoorclimate.IndoorClimateMeasurement) {
}
