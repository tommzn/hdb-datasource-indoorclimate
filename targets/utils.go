package targets

import (
	"strconv"
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

// ParseMeasurementValue will try to convert passed string value to float or integer.
// In case both convertions fail passed string value is returned.
func parseMeasurementValue(measurementValue string) interface{} {

	if strings.Contains(measurementValue, ".") {
		if floatValue, err := strconv.ParseFloat(measurementValue, 64); err == nil {
			return floatValue
		}
	} else {
		if intValue, err := strconv.Atoi(measurementValue); err == nil {
			return intValue
		}
	}
	return measurementValue
}
