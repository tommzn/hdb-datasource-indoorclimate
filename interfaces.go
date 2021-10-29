package indoorclimate

type messageTarget interface {
	send(IndorrClimate) error
}
