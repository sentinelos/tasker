// Package writers implementation of diagnostic metrics writers interface
package writers

import (
	"io"

	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

// Writer defines a logger writers interface.
type Writer interface {
	Write(io.Writer, *metrics.Set)
}
