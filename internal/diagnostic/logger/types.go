package logger

import (
	"time"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

type (
	// Severity  represents the severity of a diagnostic.
	Severity uint

	// Entry represents information to be presented to a user about a debug, info or etc. of application.
	Entry struct {
		Stamp    time.Time
		Severity Severity
		Message  string
		Labels   labels.Labels
	}
)
