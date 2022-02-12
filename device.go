package indoorclimate

import (
	"errors"
	"fmt"

	"github.com/muka/go-bluetooth/api"
)

// NewIndoorClimateSensor returns a new indoor climate device.
func NewIndoorClimateSensor(adapterId, deviceId string) SensorDevice {
	return &IndoorClimateSensor{
		adapterId: adapterId,
		deviceId:  deviceId,
	}
}

// ID returns sensor device id.
func (sensor *IndoorClimateSensor) Id() string {
	return sensor.deviceId
}

// Connect will try to connect to a device and will return with an error if failing.
func (sensor *IndoorClimateSensor) Connect() error {

	adapter, err := api.GetAdapter(sensor.adapterId)
	if err != nil {
		return err
	}

	sensor.device, err = adapter.GetDeviceByAddress(sensor.deviceId)
	if err != nil {
		return err
	}
	if sensor.device == nil {
		return fmt.Errorf("Unable to get device: %s", sensor.deviceId)
	}
	return sensor.device.Connect()
}

// Disconnect will try to disconnect from current device and returns with an error if it fails.
func (sensor *IndoorClimateSensor) Disconnect() error {
	err := sensor.device.Disconnect()
	if err == nil {
		sensor.device = nil
	}
	return err
}

// ReadValue will try to read measurment value for given characteristics.
func (sensor *IndoorClimateSensor) ReadValue(characteristicsId string) ([]byte, error) {

	if err := sensor.validateConnection(); err != nil {
		return nil, err
	}

	characteristic, err := sensor.device.GetCharByUUID(characteristicsId)
	if err != nil {
		return nil, err
	}
	return characteristic.ReadValue(nil)

}

// ValidateConnection will assert an existing, connected device.
func (sensor *IndoorClimateSensor) validateConnection() error {

	if sensor.device == nil {
		return errors.New("No device connected!")
	}
	if connected, err := sensor.device.GetConnected(); !connected || err != nil {
		return fmt.Errorf("Device %s not connected!", sensor.deviceId)
	}
	return nil
}
