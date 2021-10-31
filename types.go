package indoorclimate

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
	core "github.com/tommzn/hdb-datasource-core"
	events "github.com/tommzn/hdb-events-go"
)

// MqttClient is used to subsribe to a MQTT broker and process messages with indoor climate data.
type MqttClient struct {
	conf           config.Config
	secretsManager secrets.SecretsManager
	logger         log.Logger
	targets        []MessageTarget
}

// SqsTarget sends passed indoor climate data to a AWS SQS queue.
type SqsTarget struct {

	// publisher is a SQS client to publish messages.
	publisher core.Publisher
}

// logTarget writes given indoor climate data to an internal logger
type logTarget struct {
	logger log.Logger
}

// collectorTarget collects passed infoor climate data in local storage.
type collectorTarget struct {
	messages []events.IndoorClimate
}

// messageHandler is used to process messages received from a MQTT topic.
type messageHandler struct {
	logger log.Logger
}
