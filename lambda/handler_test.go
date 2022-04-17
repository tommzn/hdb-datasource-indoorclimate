package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	targets "github.com/tommzn/hdb-datasource-indoorclimate/targets"
)

type HandlerTestSuite struct {
	suite.Suite
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleIotEvent() {

	handler := New(loggerForTest(), loadConfigForTest(nil))
	handler.appendTarget(targets.NewStdoutTarget())
	suite.Nil(handler.HandleEvent(context.Background(), indoorClimateDataForTest()))
	suite.Nil(handler.HandleEvent(context.Background(), batteryDataForTest()))
	suite.NotNil(handler.HandleEvent(context.Background(), invalidIndoorClimateDataForTest()))
}

func (suite *HandlerTestSuite) TestBootstrap() {

	os.Setenv("AWS_REGION", "us-west-1")
	handler, err := bootstrap()
	suite.Nil(err)
	suite.NotNil(handler)
}
