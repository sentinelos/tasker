package diagnostic

import (
	"errors"
	"strings"
)

// Severity  represents the severity of a diagnostic.
type Severity uint

const (
	DiagInvalid Severity = iota
	DiagTrace
	DiagDebug
	DiagInfo
	DiagWarn
	DiagError
	DiagFatal
)

var (
	SeverityNames = map[Severity]string{
		DiagTrace: "trace",
		DiagDebug: "debug",
		DiagInfo:  "info",
		DiagWarn:  "warn",
		DiagError: "error",
		DiagFatal: "fatal",
	}

	SeverityStrings = map[string]Severity{
		"trace": DiagTrace,
		"debug": DiagDebug,
		"info":  DiagInfo,
		"warn":  DiagWarn,
		"error": DiagError,
		"fatal": DiagFatal,
	}

	// ErrDiagInvalid is returned if the severity is invalid.
	ErrDiagInvalid = errors.New("invalid diagnostic severity")
)

// String severity to string
func (s Severity) String() string {
	return SeverityNames[s]
}

// ParseSeverity parses severity string.
func ParseSeverity(str string) (Severity, error) {
	s, ok := SeverityStrings[strings.ToLower(str)]
	if !ok {
		return DiagInvalid, ErrDiagInvalid
	}

	return s, nil
}

// MustParseSeverity parses severity string or panics.
func MustParseSeverity(str string) Severity {
	s, err := ParseSeverity(str)
	if err != nil {
		panic(ErrDiagInvalid)
	}

	return s
}
