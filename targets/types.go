// Package contains provides diferent targets Indoor Climate date can be send to.
package targets

import (
	pubsub "github.com/tommzn/aws-pub-sub"
	log "github.com/tommzn/go-log"
	metrics "github.com/tommzn/go-metrics"
	core "github.com/tommzn/hdb-datasource-core"
)

// SqsTarget sends passed indoor climate data to a AWS SQS queue.
type SqsTarget struct {

	// Publisher is a SQS client to publish messages.
	publisher core.Publisher
}

// LogTarget writes given indoor climate data to an internal logger
type LogTarget struct {
	logger log.Logger
}

// TimestreamTarget writes writes publishing metrics to AWS Timestream.
type TimestreamTarget struct {
	metricPublisher metrics.Publisher
}

// StdoutTarget writes given indoor climate data to Stdout unsing fmt package.
type StdoutTarget struct {
}

// SnsTarget sends passed indoor climate data to a AWS SNS topic.
type SnsTarget struct {

	// topicArn defines target mesages should be send to.
	topicArn *string

	// Publisher is a SNS client to publish messages.
	publisher pubsub.Publisher
}
