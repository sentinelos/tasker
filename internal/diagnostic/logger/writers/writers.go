// Package writers implementation of diagnostic logger writers interface
package writers

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
)

type (
	// Writer defines a logger writers interface.
	Writer interface {
		Write(*logger.Entry)
	}
)
