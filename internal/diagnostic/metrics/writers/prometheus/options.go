package prometheus

import (
	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

func NewOptions(opt ...Option) Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// WithLabels set tags for the Prometheus writer
func WithLabels(tags labels.Labels) Option {
	return func(o *Options) {
		o.Labels = tags
	}
}
