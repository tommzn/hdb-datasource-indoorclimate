package indoorclimate

// MessageTarget is uses as destination for received indoor climate data.
type MessageTarget interface {
	Send(IndorrClimate) error
}
