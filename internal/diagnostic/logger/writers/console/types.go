package console

import (
	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

// Console defines an entry console writer.
type Console struct {
	Options
}

type Options struct {
	// ColorOutput determines if used colorized output.
	ColorOutput bool

	// QuoteString determines if quoting string values.
	QuoteString bool

	// EndWithMessage determines if output message in the end.
	EndWithMessage bool

	TimeFormat string

	Tags labels.Labels
}
