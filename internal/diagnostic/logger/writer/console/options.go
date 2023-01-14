package console

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
)

type Option func(o *ConsoleOptions)

func NewOptions(opt ...Option) ConsoleOptions {
	opts := ConsoleOptions{
		ColorOutput: true,
		QuoteString: true,
		TimeFormat:  logger.DefaultTimeStamp,
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// WithColorOutput set color output for the console writer
func WithColorOutput(colorOutput bool) Option {
	return func(o *ConsoleOptions) {
		o.ColorOutput = colorOutput
	}
}

// WithQuoteString set quote string for the console writer
func WithQuoteString(quoteString bool) Option {
	return func(o *ConsoleOptions) {
		o.QuoteString = quoteString
	}
}

// WithEndWithMessage set end with message for the console writer
func WithEndWithMessage(endWithMessage bool) Option {
	return func(o *ConsoleOptions) {
		o.EndWithMessage = endWithMessage
	}
}

// WithTimeFormat set time format for the console writer
func WithTimeFormat(timeFormat string) Option {
	return func(o *ConsoleOptions) {
		o.TimeFormat = timeFormat
	}
}
