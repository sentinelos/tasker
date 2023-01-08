package diagnostic

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logging"
	"github.com/sentinelos/tasker/internal/diagnostic/logging/writer"
)

// Diagnostic is a list of Diagnostic instances.
type Diagnostic struct {
	DiagnosticOptions
}

type DiagnosticOptions struct {
	Name        string
	Description string
	Severity    logging.Severity
	LogWriters  []writer.Writer
	Meta        map[string]string
}
