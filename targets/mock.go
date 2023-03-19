package targets

import indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"

func NewMockTarget() *MockTarget {
	return &MockTarget{
		Measurements: []indoorclimate.IndoorClimateMeasurement{},
	}
}

func (target *MockTarget) SendMeasurement(measurement indoorclimate.IndoorClimateMeasurement) error {
	target.Measurements = append(target.Measurements, measurement)
	return nil
}
