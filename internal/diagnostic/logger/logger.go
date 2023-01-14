// Package logger implementation of diagnostic logs
package logger

import (
	"time"
)

var (
	DefaultSeverity  = Info
	DefaultTimeStamp = time.RFC3339
)
