package prometheus

import (
	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

type Option func(o *Options)

func NewOptions(opt ...Option) Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// WithTags set tags for the Prometheus writer
func WithTags(tags labels.Labels) Option {
	return func(o *Options) {
		o.Tags = tags
	}
}
