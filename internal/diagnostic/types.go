package diagnostic

import (
	"time"
)

// Diagnostic is a list of Diagnostic instances.
type Diagnostic struct {
	DiagnosticOptions
}

type DiagnosticOptions struct {
	Name        string
	Description string
	Severity    Severity
	Writers     []Writer
	Meta        map[string]string
}

// Entry represents information to be presented to a user about a debug, info or etc. of application.
type Entry struct {
	Time     time.Time
	Severity Severity
	Message  string
	Fields   []Field
}

type Field struct {
	Name        string
	Description string
	Value       string
	Sensitive   bool
}

// Writer defines an entry writer interface.
type Writer interface {
	WriteEntry(*Entry)
}

// ConsoleWriter defines an entry console writer.
type ConsoleWriter struct {
	ConsoleWriterOptions
}

type ConsoleWriterOptions struct {
	// ColorOutput determines if used colorized output.
	ColorOutput bool

	// QuoteString determines if quoting string values.
	QuoteString bool

	// EndWithMessage determines if output message in the end.
	EndWithMessage bool

	TimeFormat string
}
