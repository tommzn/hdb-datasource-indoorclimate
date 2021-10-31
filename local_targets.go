package indoorclimate

import (
	log "github.com/tommzn/go-log"
	events "github.com/tommzn/hdb-events-go"
)

// newLogTarget returns a new message target which logs passed indoor climate data.
func newLogTarget(logger log.Logger) MessageTarget {
	return &logTarget{
		logger: logger,
	}
}

// Send passed indoor climate date to a logger.
func (target *logTarget) Send(indoorClimate events.IndoorClimate) error {
	target.logger.Infof("IndoorCliemate, Device: %s, Type: %s, Value: %s",
		indoorClimate.DeviceId, indoorClimate.Type, indoorClimate.Value)
	return nil
}

// newCollectorTarget returns a new message target which collects passed indoor climate data locally.
func newCollectorTarget() MessageTarget {
	return &collectorTarget{
		messages: []events.IndoorClimate{},
	}
}

// Send will append passed indoor climate data to local storage.
func (target *collectorTarget) Send(indoorClimate events.IndoorClimate) error {
	target.messages = append(target.messages, indoorClimate)
	return nil
}
