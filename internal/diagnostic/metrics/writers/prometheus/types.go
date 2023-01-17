package prometheus

import (
	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

// Prometheus defines a Prometheus writer.
type Prometheus struct {
	Options
}

type Options struct {
	Tags labels.Labels
}
