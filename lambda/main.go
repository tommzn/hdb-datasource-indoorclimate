package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
)

func main() {

	handler := bootstrap()
	lambda.Start(handler.HandleEvent)
}

func bootstrap() MessageHandler {

	conf := loadConfig()
	secretsManager := newSecretsManager()
	logger := newLogger(conf, secretsManager)
	return New(logger, conf)
}

// loadConfig from config file.
func loadConfig() config.Config {

	configSource, err := config.NewS3ConfigSourceFromEnv()
	exitOnError(err)

	conf, err := configSource.Load()
	exitOnError(err)
	return conf
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

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}
