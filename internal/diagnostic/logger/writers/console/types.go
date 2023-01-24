package console

import (
	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

type (
	// Console defines an entry console writer.
	Console struct {
		Options
	}

	Options struct {
		// ColorOutput determines if used colorized output.
		ColorOutput bool

		// QuoteString determines if quoting string values.
		QuoteString bool

		// EndWithMessage determines if output message in the end.
		EndWithMessage bool

		TimeFormat string

		Labels labels.Labels
	}

	Option func(o *Options)
)
