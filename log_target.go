package indoorclimate

import (
	"time"

	log "github.com/tommzn/go-log"
)

// newLogPublisher create a new log publisher.
func newLogPublisher(logger log.Logger) Publisher {
	return &LogPublisher{
		logger: logger,
	}
}

// SendMeasurement will write log message with level Info for passed indoor climate measurement.
func (logPublisher *LogPublisher) SendMeasurement(measurement IndoorClimateMeasurement) error {
	logPublisher.logger.Infof("IndoorClimate: DeviceId: %s, TimeStamp: %s, Type: %s, Value: %s",
		measurement.DeviceId, measurement.Timestamp.Format(time.RFC3339), measurement.Type, measurement.Value)
	return nil
}
