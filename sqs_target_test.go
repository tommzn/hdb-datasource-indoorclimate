package indoorclimate

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SqsTargetTestSuite struct {
	suite.Suite
}

func TestSqsTargetTestSuite(t *testing.T) {
	suite.Run(t, new(SqsTargetTestSuite))
}

func (suite *SqsTargetTestSuite) TestSendMessage() {

	skipCI(suite.T())

	indoorClimate := indoorCliamteDataForTest()
	publisher := NewSqsTarget(loadConfigForTest(nil), loggerForTest())
	suite.NotNil(publisher.Send(indoorClimate))
}
