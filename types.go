package indoorclimate

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
	events "github.com/tommzn/hdb-events-go"
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
	messages []events.IndoorClimate
}

// messageHandler is used to process messages received from a MQTT topic.
type messageHandler struct {
	logger log.Logger
}
