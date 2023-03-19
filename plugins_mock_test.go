package indoorclimate

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MockPluginTestSuite struct {
	suite.Suite
}

func TestMockPluginTestSuite(t *testing.T) {
	suite.Run(t, new(MockPluginTestSuite))
}

func (suite *MockPluginTestSuite) TestProcessMessages() {

	plugin := NewMockPlugin(loggerForTest(), nil)
	measurementChannel := make(chan IndoorClimateMeasurement, 10)

	plugin.SetMeasurementChannel(measurementChannel)

	plugin.MessageHandler(nil, newMqttMessage("topic", "temperature"))
	suite.Len(measurementChannel, 1)
	measurement01 := <-measurementChannel
	suite.Equal(MEASUREMENTTYPE_TEMPERATURE, measurement01.Type)

	plugin.MessageHandler(nil, newMqttMessage("topic", "humidity"))
	suite.Len(measurementChannel, 1)
	measurement02 := <-measurementChannel
	suite.Equal(MEASUREMENTTYPE_HUMIDITY, measurement02.Type)

	plugin.MessageHandler(nil, newMqttMessage("topic", "battery"))
	suite.Len(measurementChannel, 1)
	measurement03 := <-measurementChannel
	suite.Equal(MEASUREMENTTYPE_BATTERY, measurement03.Type)

	plugin.MessageHandler(nil, newMqttMessage("topic", "xxx"))
	suite.Len(measurementChannel, 1)
	measurement04 := <-measurementChannel
	suite.Equal(MEASUREMENTTYPE_TEMPERATURE, measurement04.Type)
}
