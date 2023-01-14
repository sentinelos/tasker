package logger

import (
	"time"
)

// Entry represents information to be presented to a user about a debug, info or etc. of application.
type Entry struct {
	Time     time.Time
	Severity Severity
	Message  string
	Labels   []Label
}

type Label struct {
	Name        string
	Description string
	Value       string
	Sensitive   bool
}
