package indoorclimate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
)

type LivenessestSuite struct {
	suite.Suite
}

func TestMLivenessestSuite(t *testing.T) {
	suite.Run(t, new(LivenessestSuite))
}

func (suite *LivenessestSuite) TestLivenessProve() {

	logger := loggerForTest()
	conf := loadConfigForTest(config.AsStringPtr("fixtures/testconfig_05.yml"))
	observer := NewMqttLivenessObserver(conf, logger)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	observer.Run(ctx)

	<-ctx.Done()
}
