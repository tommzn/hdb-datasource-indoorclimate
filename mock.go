package indoorclimate

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

type mqttMessage struct {
	topic, payload string
	ack            bool
}

func (message *mqttMessage) Duplicate() bool {
	return false
}

func (message *mqttMessage) Qos() byte {
	return 1
}
func (message *mqttMessage) Retained() bool {
	return false
}

func (message *mqttMessage) Topic() string {
	return message.topic
}

func (message *mqttMessage) MessageID() uint16 {
	return 1
}

func (message *mqttMessage) Payload() []byte {
	return []byte(message.payload)
}

func (message *mqttMessage) Ack() {
	message.ack = true
}

type indoorClimateSensorMock struct {
	readDelay             *time.Duration
	shouldReturnWithError bool
	connected             bool
}

func (mock *indoorClimateSensorMock) Id() string {
	return "0af0cfcb-eef8-41c0-b00f-0a307bdc9578"
}

func (mock *indoorClimateSensorMock) Connect() error {
	mock.connected = true
	return nil
}

func (mock *indoorClimateSensorMock) Disconnect() error {
	mock.connected = false
	return nil
}

func (mock *indoorClimateSensorMock) ReadValue(id string) ([]byte, error) {

	if !mock.connected {
		return []byte{}, errors.New("Not connected!")
	}
	if mock.shouldReturnWithError {
		return []byte{}, errors.New("Unable to read device data.")
	}
	if mock.readDelay != nil {
		time.Sleep(*mock.readDelay)
	}
	switch id {
	case "00002a6e-0000-1000-8000-00805f9b34fb": // Temperature
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, int64(2260))
		return buf.Bytes(), err
	case "00002a6f-0000-1000-8000-00805f9b34fb": // Humidity
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, int64(6780))
		return buf.Bytes(), err
	case "00002a19-0000-1000-8000-00805f9b34fb": // Battery
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, int64(79))
		return buf.Bytes(), err
	default:
		return []byte{}, errors.New("Unsupported: " + id)
	}
}

func newPublisherMock() *publisherMock {
	return &publisherMock{data: []IndoorClimateMeasurement{}}
}

type publisherMock struct {
	data []IndoorClimateMeasurement
}

func (mock *publisherMock) Sendeasurement(measurement IndoorClimateMeasurement) error {
	mock.data = append(mock.data, measurement)
	return nil
}
