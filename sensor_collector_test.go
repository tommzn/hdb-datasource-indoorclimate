package indoorclimate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	core "github.com/tommzn/hdb-datasource-core"
)

type SensorCollectorTestSuite struct {
	suite.Suite
}

func TestSensorCollectorTestSuite(t *testing.T) {
	suite.Run(t, new(SensorCollectorTestSuite))
}

func (suite *SensorCollectorTestSuite) TestGetSensorData() {

	publisherMock := newPublisherMock()
	collector := sensorDataCollectorForTest(publisherMock)

	collector.Run(context.Background())
	suite.Len(publisherMock.data, 3)
}

func (suite *SensorCollectorTestSuite) TestCancelRun() {

	publisherMock := newPublisherMock()
	collector := sensorDataCollectorForTest(publisherMock)
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	readDelay := 1 * time.Second
	collector.(*SensorDataCollector).devices[0].(*indoorClimateSensorMock).readDelay = &readDelay

	err := collector.Run(ctx)
	suite.NotNil(err)
	suite.Len(publisherMock.data, 0)
}

func (suite *SensorCollectorTestSuite) TestSensorDataReadErrror() {

	publisherMock := newPublisherMock()
	collector := sensorDataCollectorForTest(publisherMock)
	collector.(*SensorDataCollector).devices[0].(*indoorClimateSensorMock).shouldReturnWithError = true

	err := collector.Run(context.Background())
	suite.NotNil(err)
	suite.Len(publisherMock.data, 0)
}

func sensorDataCollectorForTest(publisher Publisher) core.Collector {

	conf := loadConfigForTest(nil)
	devices := []SensorDevice{&indoorClimateSensorMock{connected: false}}
	characteristics := characteristicsFromConfig(conf)
	return &SensorDataCollector{
		logger:          loggerForTest(),
		devices:         devices,
		characteristics: characteristics,
		publisher:       []Publisher{publisher},
		retryCount:      3,
	}
}
