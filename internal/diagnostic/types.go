package diagnostic

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
	loggerWriters "github.com/sentinelos/tasker/internal/diagnostic/logger/writers"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
	metricsWriters "github.com/sentinelos/tasker/internal/diagnostic/metrics/writers"
)

type (
	// Diagnostic is a list of Diagnostic instances.
	Diagnostic struct {
		Options
	}

	Options struct {
		Name           string
		Description    string
		Severity       logger.Severity
		LoggerWriters  map[string]loggerWriters.Writer
		MetricsWriters map[string]metricsWriters.Writer
		MetricsSet     *metrics.Set
	}

	Option func(o *Options)
)
