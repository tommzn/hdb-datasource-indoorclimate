package indoorclimate

// MockTarget can be used for testing. It appends all measurements to an internal slice.
type MockTarget struct {
	Measurements []IndoorClimateMeasurement
}

func NewMockTarget() *MockTarget {
	return &MockTarget{
		Measurements: []IndoorClimateMeasurement{},
	}
}

func (target *MockTarget) SendMeasurement(measurement IndoorClimateMeasurement) error {
	target.Measurements = append(target.Measurements, measurement)
	return nil
}
