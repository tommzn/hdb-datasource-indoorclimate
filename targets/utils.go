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
		Type:      toEventType(measurement.Type),
		Value:     measurement.Value,
	}
}

// toEventType returns corresponding measurement type from evnets package.
func toEventType(measurementType indoorclimate.MeasurementType) events.MeasurementType {

	switch measurementType {
	case indoorclimate.MEASUREMENTTYPE_TEMPERATURE:
		return events.MeasurementType_TEMPERATURE
	case indoorclimate.MEASUREMENTTYPE_HUMIDITY:
		return events.MeasurementType_HUMIDITY
	case indoorclimate.MEASUREMENTTYPE_BATTERY:
		return events.MeasurementType_BATTERY
	default:
		return events.MeasurementType_TEMPERATURE
	}
}
