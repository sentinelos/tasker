package diagnostic

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
	loggerWriters "github.com/sentinelos/tasker/internal/diagnostic/logger/writers"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
	metricsWriters "github.com/sentinelos/tasker/internal/diagnostic/metrics/writers"
)

// Diagnostic is a list of Diagnostic instances.
type Diagnostic struct {
	Options
}

type Options struct {
	Name           string
	Description    string
	Severity       logger.Severity
	LoggerWriters  map[string]loggerWriters.Writer
	MetricsWriters map[string]metricsWriters.Writer
	MetricsSet     *metrics.Set
}
