package diagnostic

import (
	"regexp"

	"github.com/sentinelos/tasker/internal/diagnostic/logging"
	"github.com/sentinelos/tasker/internal/diagnostic/logging/writer"
	"github.com/sentinelos/tasker/internal/diagnostic/logging/writer/console"
)

type Option func(o *DiagnosticOptions)

var (
	NameRegex = regexp.MustCompile(`^[a-zA-Z_.][a-zA-Z0-9_.]*$`)
)

func NewOptions(opt ...Option) DiagnosticOptions {
	opts := DiagnosticOptions{
		Severity: logging.DefaultSeverity,
		LogWriters: []writer.Writer{
			console.NewConsole(console.NewOptions()),
		},
		Meta: hostMetadata(),
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// WithName set name for the diagnostic
func WithName(name string) Option {
	return func(o *DiagnosticOptions) {
		o.Name = name
	}
}

// WithDescription set description for the diagnostic
func WithDescription(description string) Option {
	return func(o *DiagnosticOptions) {
		o.Description = description
	}
}

// WithSeverity set severity for the diagnostic
func WithSeverity(severity logging.Severity) Option {
	return func(o *DiagnosticOptions) {
		o.Severity = severity
	}
}

// WithWriter set writer for the diagnostic
func WithWriter(writer writer.Writer) Option {
	return func(o *DiagnosticOptions) {
		o.LogWriters = append(o.LogWriters, writer)
	}
}

// WithWriters set writers for the diagnostic
func WithWriters(writers []writer.Writer) Option {
	return func(o *DiagnosticOptions) {
		o.LogWriters = append(o.LogWriters, writers...)
	}
}
