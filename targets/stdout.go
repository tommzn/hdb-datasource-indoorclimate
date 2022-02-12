package targets

import (
	"fmt"
	"time"

	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

// NewStdoutTarget create a new publisher writing to stdout using fmt.
func NewStdoutTarget() indoorclimate.Publisher {
	return &StdoutTarget{}
}

// SendMeasurement will write passed indoor climate data to stdout.
func (logPublisher *StdoutTarget) SendMeasurement(measurement indoorclimate.IndoorClimateMeasurement) error {
	fmt.Printf("IndoorClimate: DeviceId: %s, TimeStamp: %s, Type: %s, Value: %s",
		measurement.DeviceId, measurement.Timestamp.Format(time.RFC3339), measurement.Type, measurement.Value)
	return nil
}
