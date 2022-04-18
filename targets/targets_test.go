package targets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
	events "github.com/tommzn/hdb-events-go"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (suite *UtilsTestSuite) TestConvertToEvent() {

	measurement := indoorclimate.IndoorClimateMeasurement{
		DeviceId:  "Devide01",
		Timestamp: time.Now(),
		Type:      indoorclimate.MEASUREMENTTYPE_TEMPERATURE,
		Value:     "21.5",
	}
	event := toIndoorClimateDate(measurement)
	suite.Equal(events.MeasurementType_TEMPERATURE, event.Type)
	suite.Equal(measurement.Value, event.Value)
	suite.Equal(measurement.DeviceId, event.DeviceId)
}

func (suite *UtilsTestSuite) TestParseMeasurementValue() {

	val1 := parseMeasurementValue("4.5")
	floatVal, ok := val1.(float64)
	suite.True(ok)
	suite.Equal(4.5, floatVal)

	val2 := parseMeasurementValue("82")
	intVal, ok := val2.(int)
	suite.True(ok)
	suite.Equal(82, intVal)

	val3 := parseMeasurementValue("xxx")
	stringVal, ok := val3.(string)
	suite.True(ok)
	suite.Equal("xxx", stringVal)
}
