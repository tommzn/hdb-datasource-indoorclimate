package indoorclimate

import (
	"strings"

	events "github.com/tommzn/hdb-events-go"
)

// toEventType returns corresponding measurement type from evnets package.
func (m MeasurementType) toEventType() events.MeasurementType {
	if val, ok := events.MeasurementType_value[strings.ToUpper(string(m))]; ok {
		return events.MeasurementType(val)
	}
	return events.MeasurementType_TEMPERATURE
}
