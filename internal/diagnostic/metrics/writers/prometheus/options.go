package prometheus

type Option func(o *Options)

func NewOptions(opt ...Option) Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// WithTags set tags for the Prometheus writer
func WithTags(tags map[string]string) Option {
	return func(o *Options) {
		o.Tags = tags
	}
}
