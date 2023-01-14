package diagnostic

import (
	"regexp"

	"github.com/sentinelos/tasker/internal/diagnostic/logger"
	"github.com/sentinelos/tasker/internal/diagnostic/logger/writer"
	"github.com/sentinelos/tasker/internal/diagnostic/logger/writer/console"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

type Option func(o *DiagnosticOptions)

var (
	NameRegex = regexp.MustCompile(`^[a-zA-Z_.][a-zA-Z0-9_.]*$`)
)

func NewOptions(opt ...Option) DiagnosticOptions {
	opts := DiagnosticOptions{
		Severity: logger.DefaultSeverity,
		LogWriters: []writer.Writer{
			console.NewConsole(console.NewOptions()),
		},
		Meta:       metadata(),
		MetricsSet: metrics.NewSet(),
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
func WithSeverity(severity logger.Severity) Option {
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
