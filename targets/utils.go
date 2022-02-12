package targets

import (
	"strings"

	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
	events "github.com/tommzn/hdb-events-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toIndoorClimateDate converts passed indoor climate measurement into an event.
func toIndoorClimateDate(measurement indoorclimate.IndoorClimateMeasurement) events.IndoorClimate {
	return events.IndoorClimate{
		DeviceId:  measurement.DeviceId,
		Timestamp: timestamppb.New(measurement.Timestamp),
		Type:      toEventType(measurement.Type),
		Value:     measurement.Value,
	}
}

// toEventType returns corresponding measurement type from evnets package.
func toEventType(measurementType indoorclimate.MeasurementType) events.MeasurementType {
	if val, ok := events.MeasurementType_value[strings.ToUpper(string(measurementType))]; ok {
		return events.MeasurementType(val)
	}
	return events.MeasurementType_TEMPERATURE
}
