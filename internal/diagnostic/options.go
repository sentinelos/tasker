package diagnostic

import (
	"github.com/sentinelos/tasker/internal/constants"
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
	loggerWriters "github.com/sentinelos/tasker/internal/diagnostic/logger/writers"
	"github.com/sentinelos/tasker/internal/diagnostic/logger/writers/console"
	"github.com/sentinelos/tasker/internal/diagnostic/metadata"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
	metricsWriters "github.com/sentinelos/tasker/internal/diagnostic/metrics/writers"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics/writers/prometheus"
)

func NewOptions(opt ...Option) Options {
	labels := metadata.Get("host").Labels
	opts := Options{
		Severity:   constants.DefaultLoggerSeverity,
		MetricsSet: metrics.NewSet(),
	}

	for _, o := range opt {
		o(&opts)
	}

	labels["diagnostic"] = opts.Name

	opts.LoggerWriters = map[string]loggerWriters.Writer{
		"console": console.NewConsole(console.NewOptions(
			console.WithLabels(labels),
		)),
	}

	opts.MetricsWriters = map[string]metricsWriters.Writer{
		"prometheus": prometheus.NewPrometheus(prometheus.NewOptions(
			prometheus.WithLabels(labels),
		)),
	}

	return opts
}

// WithName set name for the diagnostic
func WithName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// WithDescription set description for the diagnostic
func WithDescription(description string) Option {
	return func(o *Options) {
		o.Description = description
	}
}

// WithSeverity set severity for the diagnostic
func WithSeverity(severity logger.Severity) Option {
	return func(o *Options) {
		o.Severity = severity
	}
}

// WithLoggerWriter set logger writer for the diagnostic
func WithLoggerWriter(name string, writer loggerWriters.Writer) Option {
	return func(o *Options) {
		o.LoggerWriters[name] = writer
	}
}

// WithLoggerWriters set logger writer for the diagnostic
func WithLoggerWriters(writers map[string]loggerWriters.Writer) Option {
	return func(o *Options) {
		o.LoggerWriters = writers
	}
}

// WithMetricsWriter set metrics writer for the diagnostic
func WithMetricsWriter(name string, writer metricsWriters.Writer) Option {
	return func(o *Options) {
		o.MetricsWriters[name] = writer
	}
}

// WithMetricsWriters set logger writer for the diagnostic
func WithMetricsWriters(writers map[string]metricsWriters.Writer) Option {
	return func(o *Options) {
		o.MetricsWriters = writers
	}
}
