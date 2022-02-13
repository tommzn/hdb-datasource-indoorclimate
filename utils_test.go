package indoorclimate

import (
	"regexp"
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

func (suite *UtilsTestSuite) TestGenerateRandomSuffix() {

	length := 5
	randomSuffix := randStringBytes(length)
	suite.Len(randomSuffix, length)
	match, _ := regexp.MatchString("[a-z0-9]{5}", randomSuffix)
	suite.True(match)
}

func (suite *UtilsTestSuite) TestGetDeviceIdsFromConfig() {

	conf := loadConfigForTest(nil)
	devices := deviceIdsFromConfig(conf)
	suite.Len(devices, 2)
	suite.NotEqual("", devices[0])
}

func (suite *UtilsTestSuite) TestGetCharacteristicsFromConfig() {

	conf := loadConfigForTest(nil)
	characteristics := characteristicsFromConfig(conf)
	suite.Len(characteristics, 3)
	suite.NotEqual("", characteristics[0].uuid)
	suite.NotEqual("", characteristics[0].measurementType)
}
