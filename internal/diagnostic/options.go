package diagnostic

import (
	"time"
)

type Option func(o *DiagnosticOptions)

var (
	DefaultSeverity  = DiagInfo
	DefaultTimeStamp = time.RFC3339
)

func NewOptions(opt ...Option) DiagnosticOptions {
	opts := DiagnosticOptions{
		Severity: DefaultSeverity,
		Writers: []Writer{
			NewConsoleWriter(ConsoleWriterOptions{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
				TimeFormat:     DefaultTimeStamp,
			}),
		},
		Meta: HostMetadata(),
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
func WithSeverity(severity Severity) Option {
	return func(o *DiagnosticOptions) {
		o.Severity = severity
	}
}

// WithWriter set writer for the diagnostic
func WithWriter(writer Writer) Option {
	return func(o *DiagnosticOptions) {
		o.Writers = append(o.Writers, writer)
	}
}

// WithWriters set writers for the diagnostic
func WithWriters(writers []Writer) Option {
	return func(o *DiagnosticOptions) {
		o.Writers = append(o.Writers, writers...)
	}
}
