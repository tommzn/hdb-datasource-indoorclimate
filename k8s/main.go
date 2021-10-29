package main

import (
	"context"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"

	core "github.com/tommzn/hdb-datasource-core"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

func main() {

	ctx := context.Background()
	collector, err := bootstrap()
	if err != nil {
		panic(err)
	}
	collector.Run(ctx)
}

// bootstrap loads config and creates a new scheduled collector with a exchangerate datasource.
func bootstrap() (core.Collector, error) {

	conf := loadConfig()
	secretsManager := newSecretsManager(conf)
	logger := newLogger(conf, secretsManager)
	datasource := indoorclimate.New(conf, logger, secretsManager)
	return core.NewContinuousCollector(datasource, logger), nil
}

// loadConfig from config file.
func loadConfig() config.Config {

	configSource, err := config.NewS3ConfigSourceFromEnv()
	if err != nil {
		panic(err)
	}

	conf, err := configSource.Load()
	if err != nil {
		panic(err)
	}
	return conf
}

// newSecretsManager retruns a new secrets manager from passed config.
func newSecretsManager(conf config.Config) secrets.SecretsManager {
	return secrets.NewSecretsManagerByConfig(conf)
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager) log.Logger {
	return log.NewLoggerFromConfig(conf, secretsMenager)
}
