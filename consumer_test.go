package indoorclimate

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	core "github.com/tommzn/hdb-datasource-core"
)

type ConsumerTestSuite struct {
	suite.Suite
	waitGroup  *sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (suite *ConsumerTestSuite) SetupTest() {
	suite.waitGroup = &sync.WaitGroup{}
	suite.ctx, suite.cancelFunc = context.WithCancel(context.Background())
}

func TestConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}

func (suite *ConsumerTestSuite) TestHandleMessage() {

	consumer := messageConsumerForTest(nil)

	topic1 := "iobroker/ble/0/a4:f3:e6:b8:d1:c6/temperature"
	message1 := messageWithTopicForTest(topic1)
	consumer.processMessage(nil, message1)

	topic2 := "iobroker/ble/0/a4:f3:e6:b8:xx:yy/temperature"
	message2 := messageWithTopicForTest(topic2)
	consumer.processMessage(nil, message2)

	topic3 := "iobroker/ble/0/a4:f3:e6:b8:d1:c6/"
	message3 := messageWithTopicForTest(topic3)
	consumer.processMessage(nil, message3)

}

func (suite *ConsumerTestSuite) TestCredentialsProvider() {

	consumer := messageConsumerForTest(nil)

	username1, password1 := consumer.credentialsProvider()
	suite.Equal("", username1)
	suite.Equal("", password1)

	username := "test-user"
	os.Setenv("TSL_MQTT_USER", username)
	username2, password2 := consumer.credentialsProvider()
	suite.Equal("", username2)
	suite.Equal("", password2)

	password := "xyz"
	os.Setenv("TSL_MQTT_PASSWORD", password)
	username3, password3 := consumer.credentialsProvider()
	suite.Equal(username, username3)
	suite.Equal(password, password3)

	os.Unsetenv("TSL_MQTT_USER")
	os.Unsetenv("TSL_MQTT_PASSWORD")
}

func (suite *ConsumerTestSuite) TestGetMqttOptions() {

	consumer := messageConsumerForTest(nil)

	mqttOptions := consumer.mqttOptions()
	suite.Len(mqttOptions.Servers, 1)
}

func (suite *ConsumerTestSuite) TestGetTopicFilters() {

	consumer := messageConsumerForTest(nil)

	filters := consumer.mqttTopicFilters()
	suite.Len(filters, 3)
}

func (suite *ConsumerTestSuite) TestGetTopicFiltersWithPrefix() {

	consumer := messageConsumerForTest(nil)
	conf := loadConfigForTest(config.AsStringPtr("fixtures/testconfig_01.yml"))
	consumer.conf = conf

	filters := consumer.mqttTopicFilters()
	suite.Len(filters, 3)
	for topic, _ := range filters {
		suite.True(strings.HasPrefix(topic, "abc"))
	}
}

func (suite *ConsumerTestSuite) TestIntegration() {

	consumer := messageConsumerForTest(nil)
	collector := newCollectorTarget()
	consumer.targets = append(consumer.targets, collector)

	suite.runConsumer(consumer, false)
	time.Sleep(1 * time.Second)
	suite.Nil(publishTestMessage(consumer.mqttOptions()))

	time.Sleep(1 * time.Second)
	suite.cancelFunc()
	suite.waitGroup.Wait()
	suite.Len(collector.(*collectorTarget).messages, 1)
}

func (suite *ConsumerTestSuite) TestConnectionLostHandler() {

	consumer := messageConsumerForTest(nil)
	suite.runConsumer(consumer, false)

	opts := consumer.mqttOptions()
	opts.AutoReconnect = false
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	time.Sleep(100 * time.Millisecond)
	suite.cancelFunc()
	suite.waitGroup.Wait()
}

func (suite *ConsumerTestSuite) TestConnectionError() {

	consumer := messageConsumerForTest(config.AsStringPtr("fixtures/testconfig_02.yml"))
	suite.runConsumer(consumer, true)

	time.Sleep(100 * time.Millisecond)
	suite.cancelFunc()
	suite.waitGroup.Wait()
}

func (suite *ConsumerTestSuite) runConsumer(consumer core.Collector, shouldReturnWithError bool) {

	suite.waitGroup.Add(1)
	go func() {
		defer suite.waitGroup.Done()
		err := consumer.Run(suite.ctx)
		suite.Equal(shouldReturnWithError, err != nil)
	}()
}

func publishTestMessage(opts *mqtt.ClientOptions) error {

	topic := "/ble/0/1a:1a:1a:1a:1a:1a/temperature"
	opts.SetClientID("test_publisher")
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	token := mqttClient.Publish(topic, 0, false, "23.5")
	<-token.Done()
	err := token.Error()
	mqttClient.Disconnect(500)
	return err
}
