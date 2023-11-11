package targets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	config "github.com/tommzn/go-config"
	utils "github.com/tommzn/go-utils"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

// NewS3Target create a new S3 uploader.
func NewS3Target(conf config.Config) (indoorclimate.Publisher, error) {

	bucket := conf.Get("aws.s3.bucket", nil)
	if bucket == nil {
		return nil, errors.New("No S3 bucket defined.")
	}

	return &S3Target{
		awsConfig: newAWSConfig(conf),
		bucket:    bucket,
		path:      conf.Get("aws.s3.path", nil),
	}, nil
}

// newAWSConfig try to find AWS region in passed config or in environment variable AWS_REGION
// and returns a new AWS config.
func newAWSConfig(conf config.Config) *aws.Config {

	awsConfig := aws.NewConfig()

	if conf != nil {
		configKeys := []string{"aws.s3.region", "aws.region"}
		for _, configKey := range configKeys {
			if awsRegion := conf.Get(configKey, nil); awsRegion != nil {
				return awsConfig.WithRegion(*awsRegion)
			}
		}

	}

	if awsRegion, ok := os.LookupEnv("AWS_REGION"); ok {
		return awsConfig.WithRegion(awsRegion)
	}

	return awsConfig
}

// SendMeasurement will start to transfer passed measurement to a target.
func (target *S3Target) SendMeasurement(measurement indoorclimate.IndoorClimateMeasurement) error {

	errStack := utils.NewErrorStack()

	jsonData, jsonErr := json.Marshal(measurement)
	errStack.Append(jsonErr)

	_, s3Err := target.getS3Uploader().Upload(
		&s3manager.UploadInput{
			Bucket: target.bucket,
			Key:    target.getObjectKey(measurement),
			Body:   bytes.NewReader(jsonData),
		})
	errStack.Append(s3Err)
	return errStack.AsError()
}

// sS3Session creates a new AWS session, if none exists.
func (target *S3Target) s3Session() *session.Session {
	if target.awsSession == nil {
		target.awsSession = session.Must(session.NewSession(target.awsConfig))
	}
	return target.awsSession
}

// GetS3Uploader returns current uploader. If none exists a new will be created.
func (target *S3Target) getS3Uploader() *s3manager.Uploader {
	if target.s3Uploader == nil {
		target.s3Uploader = s3manager.NewUploader(target.s3Session())
	}
	return target.s3Uploader
}

// GetObjectKey generates a S3 object key for given measurement.
func (target *S3Target) getObjectKey(measurement indoorclimate.IndoorClimateMeasurement) *string {
	objectKey := fmt.Sprintf("%s/%s.json", measurement.DeviceId, measurement.Type)
	if target.path != nil {
		objectKey = fmt.Sprintf("%s/%s", *target.path, objectKey)
	}
	return &objectKey
}
