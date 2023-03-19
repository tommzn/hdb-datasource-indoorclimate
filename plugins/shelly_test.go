package plugins

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

type ShellyHTPluginTestSuite struct {
	suite.Suite
}

func TestShellyHTPluginTestSuite(t *testing.T) {
	suite.Run(t, new(ShellyHTPluginTestSuite))
}

func (suite *ShellyHTPluginTestSuite) TestProcessMessage() {

	measurementChannel := make(chan indoorclimate.IndoorClimateMeasurement, 10)
	messagePayload, err := os.ReadFile("fixtures/shelly_message_01.json")
	suite.Nil(err)

	plugin := NewShellyHTPlugin(loggerForTest())
	plugin.SetMeasurementChannel(measurementChannel)

	plugin.MessageHandler(nil, newMqttMessage("topic", string(messagePayload)))
	suite.Len(measurementChannel, 3)
	measurement01 := <-measurementChannel
	suite.Equal(indoorclimate.MEASUREMENTTYPE_TEMPERATURE, measurement01.Type)
}

func (suite *ShellyHTPluginTestSuite) TestProcessOmitMessage() {

	measurementChannel := make(chan indoorclimate.IndoorClimateMeasurement, 10)
	messagePayload, err := os.ReadFile("fixtures/shelly_message_02.json")
	suite.Nil(err)

	plugin := NewShellyHTPlugin(loggerForTest())
	plugin.SetMeasurementChannel(measurementChannel)

	plugin.MessageHandler(nil, newMqttMessage("topic", string(messagePayload)))
	suite.Len(measurementChannel, 0)
}
