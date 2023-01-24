package prometheus

import (
	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

// Prometheus defines a Prometheus writer.
type (
	Prometheus struct {
		Options
	}

	Options struct {
		Labels labels.Labels
	}

	Option func(o *Options)
)
