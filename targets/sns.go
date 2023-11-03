package targets

import (
	"errors"

	pubsub "github.com/tommzn/aws-pub-sub"
	config "github.com/tommzn/go-config"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

// NewSnsTarget creates a new publisher for AWS SNS.
func NewSnsTarget(conf config.Config) indoorclimate.Publisher {
	topicArn := conf.Get("hdb.topic.arn", nil)
	return &SnsTarget{
		topicArn:  topicArn,
		publisher: pubsub.NewSnsPublisher(conf),
	}
}

// SendMeasurement will start to transfer passed measurement to a target.
func (target *SnsTarget) SendMeasurement(measurement indoorclimate.IndoorClimateMeasurement) error {

	if target.topicArn == nil {
		return errors.New("No target topic specified.")
	}

	event := toIndoorClimateDate(measurement)
	return target.publisher.Send(*target.topicArn, &event)
}
