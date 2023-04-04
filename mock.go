package indoorclimate

func newPublisherMock() *publisherMock {
	return &publisherMock{data: []IndoorClimateMeasurement{}}
}

type publisherMock struct {
	data []IndoorClimateMeasurement
}

func (mock *publisherMock) SendMeasurement(measurement IndoorClimateMeasurement) error {
	mock.data = append(mock.data, measurement)
	return nil
}
