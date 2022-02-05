package indoorclimate

import (
	core "github.com/tommzn/hdb-datasource-core"
	events "github.com/tommzn/hdb-events-go"
)

// Collector will fetch data from a datasource a process it.
type Collector interface {

	// Core collector interface.
	// See https://github.com/tommzn/hdb-datasource-core/blob/main/interfaces.go
	core.Collector

	// AppendMessageTarget adds passed message target to internal stack.
	AppendMessageTarget(MessageTarget)
}

// MessageTarget is uses as destination for received indoor climate data.
type MessageTarget interface {

	// Send passed indoor climate data to defined destination.
	Send(events.IndoorClimate) error
}

// Publisher sends given measuremnts to different targets.
type Publisher interface {

	// Send will start to transfer passed measurement to a target.
	Send(IndoorClimateMeasurement) error
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
