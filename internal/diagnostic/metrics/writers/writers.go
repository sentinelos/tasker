// Package writers implementation of diagnostic metrics writers interface
package writers

import (
	"io"

	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

type (
	// Writer defines a logger writers interface.
	Writer interface {
		Write(io.Writer, *metrics.Set)
	}
)
