package indoorclimate

import (
	"context"
	"fmt"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

type MqttCollectorTestSuite struct {
	suite.Suite
	mqttClient mqtt.Client
	topic      string
	logger     log.Logger
}

func TestMqttCollectorTestSuite(t *testing.T) {
	suite.Run(t, new(MqttCollectorTestSuite))
}

func (suite *MqttCollectorTestSuite) SetupTest() {
	suite.logger = loggerForTest()
	suite.topic = fmt.Sprintf("testtopic%d", time.Now().Unix())
}

func (suite *MqttCollectorTestSuite) connectToMqttBroker() {
	options := mqtt.NewClientOptions()
	options.SetOrderMatters(false)
	options.AddBroker("tcp://localhost:1883")
	options.SetClientID(fmt.Sprintf("client_%d", time.Now().Unix()))
	client := mqtt.NewClient(options)
	token := client.Connect()
	token.Wait()
	suite.Nil(token.Error())
	suite.mqttClient = client
}

func (suite *MqttCollectorTestSuite) TearDownTest() {
	if suite.mqttClient != nil {
		suite.mqttClient.Disconnect(250)
	}
}

func (suite *MqttCollectorTestSuite) TestCreateClient() {

	conf := loadConfigForTest(config.AsStringPtr("fixtures/testconfig_03.yml"))

	collector01 := NewMqttCollector(conf, suite.logger, secretsManagerForTest())
	suite.NotNil(collector01)
	suite.Len(collector01.publisher, 0)
	suite.Len(collector01.subscriptions, 0)
}

func (suite *MqttCollectorTestSuite) TestConsumeIndoorClimateData() {

	// Setuo
	conf := loadConfigForTest(config.AsStringPtr("fixtures/testconfig_03.yml"))
	collector := NewMqttCollector(conf, suite.logger, secretsManagerForTest())
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Prepare
	suite.connectToMqttBroker()
	suite.subscribeToTestTopic(collector)
	mockTarget := NewMockTarget()
	collector.AppendTarget(mockTarget)

	// Run
	go collector.Run(ctx)
	time.Sleep(3 * time.Second)
	suite.sendMessage("humidity")
	suite.sendMessage("temperature")
	suite.sendMessage("battery")
	suite.sendMessage("xxx")

	// Stop
	time.Sleep(3 * time.Second)
	cancelFunc()
	time.Sleep(1 * time.Second)

	// Assertions
	suite.Len(collector.measurements, 0)
	suite.Len(mockTarget.Measurements, 4)
}

func (suite *MqttCollectorTestSuite) TestBrokerConnectError() {

	conf := loadConfigForTest(config.AsStringPtr("fixtures/testconfig_04.yml"))
	collector := NewMqttCollector(conf, suite.logger, secretsManagerForTest())
	ctx, cancelFunc := context.WithCancel(context.Background())

	suite.NotNil(collector.Run(ctx))
	cancelFunc()
}

func (suite *MqttCollectorTestSuite) subscribeToTestTopic(collector *MqttCollector) {
	plugin := NewMockPlugin(suite.logger, nil)
	collector.subscriptions = append(collector.subscriptions, MqttSubscriptionConfig{Topic: suite.topic, Plugin: plugin})
}

func (suite *MqttCollectorTestSuite) sendMessage(message string) {
	token := suite.mqttClient.Publish(suite.topic, 1, true, message)
	token.Wait()
	suite.Nil(token.Error())
}
