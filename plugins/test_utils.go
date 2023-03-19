package plugins

import (
	log "github.com/tommzn/go-log"
)

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}
