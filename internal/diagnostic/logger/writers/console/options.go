package console

import (
	"github.com/sentinelos/tasker/internal/constants"
	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

func NewOptions(opt ...Option) Options {
	opts := Options{
		ColorOutput: true,
		QuoteString: true,
		TimeFormat:  constants.DefaultLoggerTimeStamp,
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// WithColorOutput set color output for the Console writer
func WithColorOutput(colorOutput bool) Option {
	return func(o *Options) {
		o.ColorOutput = colorOutput
	}
}

// WithQuoteString set quote string for the Console writer
func WithQuoteString(quoteString bool) Option {
	return func(o *Options) {
		o.QuoteString = quoteString
	}
}

// WithEndWithMessage set end with message for the Console writer
func WithEndWithMessage(endWithMessage bool) Option {
	return func(o *Options) {
		o.EndWithMessage = endWithMessage
	}
}

// WithTimeFormat set time format for the Console writer
func WithTimeFormat(timeFormat string) Option {
	return func(o *Options) {
		o.TimeFormat = timeFormat
	}
}

// WithLabels set tags for the Prometheus writer
func WithLabels(labels labels.Labels) Option {
	return func(o *Options) {
		o.Labels = labels
	}
}
