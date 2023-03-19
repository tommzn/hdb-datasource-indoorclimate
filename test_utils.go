package indoorclimate

import (
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

func newTestTarget() *testTarget {
	return &testTarget{
		measurements: []IndoorClimateMeasurement{},
	}
}

type testTarget struct {
	measurements []IndoorClimateMeasurement
}

func (target *testTarget) SendMeasurement(measurement IndoorClimateMeasurement) error {
	target.measurements = append(target.measurements, measurement)
	return nil
}
