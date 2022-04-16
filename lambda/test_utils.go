package main

import (
	"time"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
)

// loadConfigForTest loads test config.
func loadConfigForTest(fileName *string) config.Config {

	configFile := "fixtures/testconfig.yml"
	if fileName != nil {
		configFile = *fileName
	}
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

// secretsManagerForTest returns a default secrets manager for testing
func secretsManagerForTest() secrets.SecretsManager {
	return secrets.NewSecretsManager()
}

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

func indoorClimateDateForTest() IndoorClimateDate {
	return IndoorClimateDate{
		DeviceId:       "XYZ",
		Characteristic: "temp",
		TimeStamp:      time.Now().Unix(),
		Value:          "Fgg=",
	}
}
