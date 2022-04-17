package main

import (
	"time"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
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

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

func indoorClimateDataForTest() IndoorClimateData {
	return IndoorClimateData{
		DeviceId:       "XYZ",
		Characteristic: "temperature",
		TimeStamp:      time.Now().Unix(),
		Value:          "Fgg=",
	}
}

func invalidIndoorClimateDataForTest() IndoorClimateData {
	return IndoorClimateData{
		DeviceId:       "XYZ",
		Characteristic: "temperature",
		TimeStamp:      time.Now().Unix(),
		Value:          "+&/(",
	}
}

func batteryDataForTest() IndoorClimateData {
	return IndoorClimateData{
		DeviceId:       "XYZ",
		Characteristic: "battery",
		TimeStamp:      time.Now().Unix(),
		Value:          "Og==",
	}
}
