package logger

import (
	"errors"
	"strings"
)

// Severity  represents the severity of a diagnostic.
type Severity uint

const (
	Invalid Severity = iota
	Trace
	Debug
	Info
	Warn
	Error
	Fatal
)

var (
	SeverityNames = map[Severity]string{
		Trace: "trace",
		Debug: "debug",
		Info:  "info",
		Warn:  "warn",
		Error: "error",
		Fatal: "fatal",
	}

	SeverityStrings = map[string]Severity{
		"trace": Trace,
		"debug": Debug,
		"info":  Info,
		"warn":  Warn,
		"error": Error,
		"fatal": Fatal,
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
		return Invalid, ErrDiagInvalid
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
