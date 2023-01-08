package diagnostic

import (
	"github.com/sentinelos/tasker/internal/diagnostic/logging"
)

var (
	diagnostics = map[string]*Diagnostic{}
)

func init() {
	Register(NewDiagnostic(NewOptions(
		WithName("default"),
	)))
}

func Register(diagnostic *Diagnostic) {
	if _, ok := diagnostics[diagnostic.Name]; !ok {
		diagnostics[diagnostic.Name] = diagnostic
	}
}

func Use(diagnostic string) *Diagnostic {
	if diag, ok := diagnostics[diagnostic]; ok {
		return diag
	}

	diag := NewDiagnostic(NewOptions(
		WithName(diagnostic),
	))

	Register(diag)

	return diag
}

func Trace(diagnostic, message string, fields ...logging.Label) {
	Use(diagnostic).Trace(message, fields...)
}

func Debug(diagnostic, message string, fields ...logging.Label) {
	Use(diagnostic).Debug(message, fields...)
}

func Info(diagnostic, message string, fields ...logging.Label) {
	Use(diagnostic).Info(message, fields...)
}

func Warn(diagnostic, message string, fields ...logging.Label) {
	Use(diagnostic).Warn(message, fields...)
}

func Error(diagnostic, message string, fields ...logging.Label) {
	Use(diagnostic).Error(message, fields...)
}

func Fatal(diagnostic, message string, fields ...logging.Label) {
	Use(diagnostic).Fatal(message, fields...)
}
