package indoorclimate

// Publisher sends given measuremnts to different targets.
type Publisher interface {

	// SendMeasurement will start to transfer passed measurement to a target.
	SendMeasurement(IndoorClimateMeasurement) error
}

// SensorDevice represents a device to fetch indoor cliamte data.
type SensorDevice interface {

	// Returns the id of current sensor device.
	Id() string

	// Connect will try to connect to a device and will return with an error if failing.
	Connect() error

	// Disconnect will try to disconnect from current device and returns with an error if it fails.
	Disconnect() error

	// ReadValue will try to read measurment value for given characteristics.
	ReadValue(string) ([]byte, error)
}
