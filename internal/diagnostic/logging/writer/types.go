package writer

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logging"
)

// Writer defines a logging writer interface.
type Writer interface {
	Write(string, *logging.Entry)
}
