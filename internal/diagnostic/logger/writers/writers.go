// Package writers implementation of diagnostic logger writers interface
package writers

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
)

// Writer defines a logger writers interface.
type Writer interface {
	Write(*logger.Entry)
}
