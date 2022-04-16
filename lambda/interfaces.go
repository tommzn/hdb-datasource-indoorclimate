package main

import "context"

type MessageHandler interface {
	HandleEvent(context.Context, IndoorClimateDate) error
}
