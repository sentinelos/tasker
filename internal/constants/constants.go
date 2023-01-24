// Package constants defines common values.
package constants

import (
	"time"

	"github.com/sentinelos/tasker/internal/diagnostic/logger"
)

const (
	// AppName is the application name.
	AppName = "tasker"

	// BadIdentifierDetail A consistent detail message for all "not a valid identifier" diagnostics.
	BadIdentifierDetail = "A name must start with a letter or underscore and may contain only letters, digits, underscores, and dashes."

	// DefaultTimeout is used if there is no timeout given
	DefaultTimeout = 10 * time.Minute

	// VarEnvPrefix Prefix for collecting variable values from environment variables
	VarEnvPrefix = "TF_VAR_"

	DefaultLoggerSeverity  = logger.Info
	DefaultLoggerTimeStamp = time.RFC3339
)
