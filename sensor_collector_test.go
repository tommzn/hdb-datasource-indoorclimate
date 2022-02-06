package indoorclimate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
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
	collector := sensorDataCollectorForTest(publisherMock, nil)

	collector.Run(context.Background())
	suite.Len(publisherMock.data, 3)
}

func (suite *SensorCollectorTestSuite) TestGetSensorDataWithSchedule() {

	publisherMock := newPublisherMock()
	collector := sensorDataCollectorForTest(publisherMock, config.AsStringPtr("fixtures/testconfig02.yml"))

	ctx, cancel := context.WithCancel(context.Background())
	go collector.Run(ctx)

	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
	suite.Len(publisherMock.data, 6)
}

func (suite *SensorCollectorTestSuite) TestCancelRun() {

	publisherMock := newPublisherMock()
	collector := sensorDataCollectorForTest(publisherMock, nil)
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	readDelay := 1 * time.Second
	collector.(*SensorDataCollector).devices[0].(*indoorClimateSensorMock).readDelay = &readDelay

	err := collector.Run(ctx)
	suite.NotNil(err)
	suite.Len(publisherMock.data, 0)
}

func (suite *SensorCollectorTestSuite) TestSensorDataReadErrror() {

	publisherMock := newPublisherMock()
	collector := sensorDataCollectorForTest(publisherMock, nil)
	collector.(*SensorDataCollector).devices[0].(*indoorClimateSensorMock).shouldReturnWithError = true

	err := collector.Run(context.Background())
	suite.NotNil(err)
	suite.Len(publisherMock.data, 0)
}

func sensorDataCollectorForTest(publisher Publisher, configFile *string) core.Collector {

	conf := loadConfigForTest(configFile)
	devices := []SensorDevice{&indoorClimateSensorMock{connected: false}}
	characteristics := characteristicsFromConfig(conf)
	collector := NewSensorDataCollector(conf, loggerForTest())
	collector.(*SensorDataCollector).devices = devices
	collector.(*SensorDataCollector).characteristics = characteristics
	collector.(*SensorDataCollector).publisher = []Publisher{publisher}
	return collector
}
