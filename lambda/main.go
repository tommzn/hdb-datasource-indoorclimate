package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
	targets "github.com/tommzn/hdb-datasource-indoorclimate/targets"
)

func main() {

	handler, err := bootstrap()
	if err != nil {
		panic(err)
	}
	lambda.Start(handler.HandleEvent)
}

func bootstrap() (MessageHandler, error) {

	conf, err := loadConfig()
	if err != nil {
		return nil, err
	}
	secretsManager := newSecretsManager()
	logger := newLogger(conf, secretsManager)
	handler := New(logger, conf)
	if queue := conf.Get("hdb.queue", nil); queue != nil {
		handler.appendTarget(targets.NewSqsTarget(conf, logger))
	}
	if timestreamTable := conf.Get("aws.timestream.table", nil); timestreamTable != nil {
		handler.appendTarget(targets.NewTimestreamTarget(conf, logger))
	}
	return handler, nil
}

// loadConfig from config file.
func loadConfig() (config.Config, error) {

	configSource, _ := config.NewS3ConfigSourceFromEnv()
	if conf, err := configSource.Load(); err == nil {
		return conf, err
	}
	return config.NewFileConfigSource(nil).Load()
}

// newSecretsManager retruns a new secrets manager from passed config.
func newSecretsManager() secrets.SecretsManager {
	return secrets.NewSecretsManager()
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager) log.Logger {
	logger := log.NewLoggerFromConfig(conf, secretsMenager)
	logContextValues := make(map[string]string)
	logContextValues[log.LogCtxNamespace] = "hdb-datasource-indoorclimate"
	logger.WithContext(log.LogContextWithValues(context.Background(), logContextValues))
	return logger
}
