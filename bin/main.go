package main

import (
	"context"
	"flag"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"

	core "github.com/tommzn/hdb-datasource-core"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "configfile", "/etc/hdb/config.yml", "Full path to config file.")
}

func main() {

	ctx := context.Background()
	collector := bootstrap(ctx)

	err := collector.Run(ctx)
	exitOnError(err)
}

// bootstrap loads config and creates a new scheduled collector with a exchangerate datasource.
func bootstrap(ctx context.Context) core.Collector {

	secretsManager := newSecretsManager()
	conf := loadConfig()
	logger := newLogger(conf, secretsManager, ctx)
	datacollector := indoorclimate.NewSensorDataCollector(conf, logger)
	return core.NewContinuousCollector(datacollector, logger)
}

// loadConfig from config file.
func loadConfig() config.Config {

	conf, err := config.NewFileConfigSource(&configFile).Load()
	exitOnError(err)
	return conf
}

// newSecretsManager retruns a new container secrets manager
func newSecretsManager() secrets.SecretsManager {
	return secrets.NewSecretsManager()
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager, ctx context.Context) log.Logger {
	logger := log.NewLoggerFromConfig(conf, secretsMenager)
	logContextValues := make(map[string]string)
	logContextValues[log.LogCtxNamespace] = "hdb-datasource-indoorclimate-ble"
	logger.WithContext(log.LogContextWithValues(ctx, logContextValues))
	return logger
}

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}
