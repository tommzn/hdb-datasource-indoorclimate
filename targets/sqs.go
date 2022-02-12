package targets

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hdb-datasource-core"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

// NewSqsTarget creates a new publisher for AWS SQS.
func NewSqsTarget(conf config.Config, logger log.Logger) indoorclimate.Publisher {
	return &SqsTarget{
		publisher: core.NewPublisher(conf, logger),
	}
}

// SendMeasurement will start to transfer passed measurement to a target.
func (target *SqsTarget) SendMeasurement(measurement indoorclimate.IndoorClimateMeasurement) error {
	event := toIndoorClimateDate(measurement)
	return target.publisher.Send(&event)
}
