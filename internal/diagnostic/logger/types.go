package logger

import (
	"time"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

// Entry represents information to be presented to a user about a debug, info or etc. of application.
type Entry struct {
	Stamp    time.Time
	Severity Severity
	Message  string
	Labels   labels.Labels
}
