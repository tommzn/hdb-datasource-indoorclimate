package indoorclimate

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
)

type MqttClient struct {
	conf           config.Config
	secretsManager secrets.SecretsManager
	logger         log.Logger
	targets        []MessageTarget
}

type logTarget struct {
	logger log.Logger
}

type collectorTarget struct {
	messages []IndorrClimate
}

// messageHandler is used to process messages received from a MQTT topic.
type messageHandler struct {
	logger log.Logger
}

type IndorrClimate struct {
	DeviceId string
	Reading  Measurement
}

type Measurement struct {
	Type, Value string
}
