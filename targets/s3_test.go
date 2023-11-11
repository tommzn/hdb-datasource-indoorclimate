package targets

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

type S3TargetTestSuite struct {
	suite.Suite
	conf config.Config
}

func TestS3TargetTestSuite(t *testing.T) {
	suite.Run(t, new(S3TargetTestSuite))
}

func (suite *S3TargetTestSuite) SetupSuite() {

	suite.skipCI()

	configFile := "s3.test.config.yml"
	configLoader := config.NewFileConfigSource(&configFile)

	var err error
	suite.conf, err = configLoader.Load()
	suite.Nil(err)
}

func (suite *S3TargetTestSuite) skipCI() {
	if _, isSet := os.LookupEnv("CI"); isSet {
		suite.T().SkipNow()
	}
}

func (suite *S3TargetTestSuite) TestUploadMeasurement() {

	s3Target, err := NewS3Target(suite.conf)
	suite.Nil(err)

	measurement := indoorclimate.IndoorClimateMeasurement{
		DeviceId:  "TestDevice01",
		Timestamp: time.Now(),
		Type:      indoorclimate.MEASUREMENTTYPE_TEMPERATURE,
		Value:     "23.5",
	}
	suite.Nil(s3Target.SendMeasurement(measurement))
}
