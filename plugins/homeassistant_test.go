package plugins

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

type HomeAssistantPluginTestSuite struct {
	suite.Suite
}

func TestHomeAssistantPluginTestSuite(t *testing.T) {
	suite.Run(t, new(HomeAssistantPluginTestSuite))
}

func (suite *HomeAssistantPluginTestSuite) TestProcessMessage() {

	measurementChannel := make(chan indoorclimate.IndoorClimateMeasurement, 10)
	messagePayload, err := os.ReadFile("fixtures/homeassistant_sensor_01.json")
	suite.Nil(err)

	plugin := NewHomeAssistantPlugin(loggerForTest())
	plugin.SetMeasurementChannel(measurementChannel)

	plugin.MessageHandler(nil, newMqttMessage("topic", string(messagePayload)))
	suite.Len(measurementChannel, 1)
	measurement01 := <-measurementChannel
	suite.Equal(indoorclimate.MEASUREMENTTYPE_HUMIDITY, measurement01.Type)
}

func (suite *HomeAssistantPluginTestSuite) TestProcessOmitMessage() {

	measurementChannel := make(chan indoorclimate.IndoorClimateMeasurement, 10)
	messagePayload, err := os.ReadFile("fixtures/homeassistant_sensor_02.json")
	suite.Nil(err)

	plugin := NewHomeAssistantPlugin(loggerForTest())
	plugin.SetMeasurementChannel(measurementChannel)

	plugin.MessageHandler(nil, newMqttMessage("topic", string(messagePayload)))
	suite.Len(measurementChannel, 0)
}
