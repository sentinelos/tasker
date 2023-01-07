package diagnostic

import (
	"time"
)

type Field struct {
	Name        string
	Description string
	Value       string
	Sensitive   bool
}

// Entry represents information to be presented to a user about a debug, info or etc. of application.
type Entry struct {
	Time     time.Time
	Severity Severity
	Message  string
	Fields   []Field
}

// Diagnostic is a list of Diagnostic instances.
type Diagnostic struct {
	Name        string
	Description string
	Severity    Severity
	Writers     []Writer
	Meta        map[string]string
}

// Writer defines an entry writer interface.
type Writer interface {
	WriteEntry(*Entry)
}
