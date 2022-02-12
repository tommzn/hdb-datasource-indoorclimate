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
		Type:      events.MeasurementType_TEMPERATURE,
		Value:     "21.5",
	}
	event := toIndoorClimateDate(measurement)
	suite.Equal(measurement.Type, event.Type)
	suite.Equal(measurement.Value, event.Value)
	suite.Equal(measurement.DeviceId, event.DeviceId)
}
