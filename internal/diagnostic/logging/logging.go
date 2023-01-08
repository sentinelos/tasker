// Package logging implementation of diagnostic logs
package logging

import (
	"time"
)

var (
	DefaultSeverity  = Info
	DefaultTimeStamp = time.RFC3339
)
