package indoorclimate

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hdb-datasource-core"
	events "github.com/tommzn/hdb-events-go"
)

// NewSqsTarget creates a new publisher for AWS SQS.
func NewSqsTarget(conf config.Config, logger log.Logger) MessageTarget {
	return &SqsTarget{
		publisher: core.NewPublisher(conf, logger),
	}
}

// Send given indoor climate data to AWS SQS queue.
func (target *SqsTarget) Send(indoorClimate events.IndoorClimate) error {
	return target.publisher.Send(&indoorClimate)
}
