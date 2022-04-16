package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleIotEvent() {

	handler := New(loggerForTest(), loadConfigForTest(nil))
	suite.Nil(handler.HandleEvent(context.Background(), indoorClimateDateForTest()))
}
