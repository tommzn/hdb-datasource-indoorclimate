package indoorclimate

import (
	"context"
	"sync"
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
	wg := &sync.WaitGroup{}
	wg.Add(1)

	observer.Run(ctx, wg)

	<-ctx.Done()
	wg.Wait()
}
