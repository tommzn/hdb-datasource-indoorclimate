package indoorclimate

import log "github.com/tommzn/go-log"

func newLogTarget(logger log.Logger) messageTarget {
	return &logTarget{
		logger: logger,
	}
}

func (target *logTarget) send(indoorClimate IndorrClimate) error {
	target.logger.Infof("IndoorCliemate, Device: %s, Type: %s, Value: %s",
		indoorClimate.DeviceId, indoorClimate.Reading.Type, indoorClimate.Reading.Value)
	return nil
}

func newCollectorTarget() messageTarget {
	return &collectorTarget{
		messages: []IndorrClimate{},
	}
}

func (target *collectorTarget) send(indoorClimate IndorrClimate) error {
	target.messages = append(target.messages, indoorClimate)
	return nil
}
