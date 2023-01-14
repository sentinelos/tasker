package writer

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
)

// Writer defines a logger writer interface.
type Writer interface {
	Write(string, *logger.Entry)
}
