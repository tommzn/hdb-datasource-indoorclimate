package targets

import (
	"time"

	log "github.com/tommzn/go-log"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

// NewLogTarget create a new log publisher.
func NewLogTarget(logger log.Logger) indoorclimate.Publisher {
	return &LogTarget{
		logger: logger,
	}
}

// SendMeasurement will write log message with level Info for passed indoor climate measurement.
func (logPublisher *LogTarget) SendMeasurement(measurement indoorclimate.IndoorClimateMeasurement) error {
	logPublisher.logger.Infof("IndoorClimate: DeviceId: %s, TimeStamp: %s, Type: %s, Value: %s",
		measurement.DeviceId, measurement.Timestamp.Format(time.RFC3339), measurement.Type, measurement.Value)
	return nil
}
