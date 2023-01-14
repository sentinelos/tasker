package diagnostic

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
	"github.com/sentinelos/tasker/internal/diagnostic/logger/writer"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

// Diagnostic is a list of Diagnostic instances.
type Diagnostic struct {
	DiagnosticOptions
}

type DiagnosticOptions struct {
	Name        string
	Description string
	Severity    logger.Severity
	LogWriters  []writer.Writer
	Meta        map[string]string
	MetricsSet  *metrics.Set
}
