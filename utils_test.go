package indoorclimate

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (suite *UtilsTestSuite) TestCreateTopics() {

	prefix := "xyz"
	topics := topicsToSubsrcibe(&prefix)
	suite.Len(topics, 3)
}

func (suite *UtilsTestSuite) TestMockTestCoverage() {

	mockedMessage := messageForTest()
	suite.False(mockedMessage.Duplicate())
	suite.True(mockedMessage.Qos() > 0)
	suite.False(mockedMessage.Retained())
	suite.True(mockedMessage.MessageID() > 0)
	mockedMessage.Ack()
}

func (suite *UtilsTestSuite) TestExtractDeviceId() {

	suite.Nil(extractDeviceId("xxx"))
	suite.Nil(extractDeviceId("iobroker/ble/0/4d:3e:4f:5c:6x:xx/temperature"))

	macAddress := "a4:f3:e6:b8:d1:c6"
	deviceId := extractDeviceId("iobroker/ble/0/" + macAddress + "/temperature")
	suite.NotNil(deviceId)
	suite.Equal(macAddress, *deviceId)
}

func (suite *UtilsTestSuite) TestExtractMeasurementType() {

	suffix1 := extractMeasurementType("iobroker/ble/0/a4:f3:e6:b8:d1:c6/temperature")
	suite.NotNil(suffix1)
	suite.Equal("temperature", *suffix1)

	suite.Nil(extractMeasurementType("iobroker.ble.0.a4:f3:e6:b8:d1:c6.temperature"))
	suite.Nil(extractMeasurementType("iobroker/ble/0/a4:f3:e6:b8:d1:c6/temperature/"))
}
