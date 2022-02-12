package targets

import (
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
	events "github.com/tommzn/hdb-events-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toIndoorClimateDate converts passed indoor climate measurement into an event.
func toIndoorClimateDate(measurement indoorclimate.IndoorClimateMeasurement) events.IndoorClimate {
	return events.IndoorClimate{
		DeviceId:  measurement.DeviceId,
		Timestamp: timestamppb.New(measurement.Timestamp),
		Type:      measurement.Type,
		Value:     measurement.Value,
	}
}
