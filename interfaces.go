package indoorclimate

import core "github.com/tommzn/hdb-datasource-core"

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
	Send(IndorrClimate) error
}
