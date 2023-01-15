package diagnostic

import (
	"io"

	"github.com/sentinelos/tasker/internal/diagnostic/logger"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

var (
	diagnostic = &Diagnostic{}
)

func init() {
	diagnostic = NewDiagnostic(NewOptions(
		WithName("default"),
	))
}

func Trace(message string, fields ...logger.Label) {
	diagnostic.Trace(message, fields...)
}

func Debug(message string, fields ...logger.Label) {
	diagnostic.Debug(message, fields...)
}

func Info(message string, fields ...logger.Label) {
	diagnostic.Info(message, fields...)
}

func Warn(message string, fields ...logger.Label) {
	diagnostic.Warn(message, fields...)
}

func Error(message string, fields ...logger.Label) {
	diagnostic.Error(message, fields...)
}

func Fatal(message string, fields ...logger.Label) {
	diagnostic.Fatal(message, fields...)
}

func Counter(name, description string) *metrics.Counter {
	return diagnostic.Counter(name, description)
}

func Gauge(name, description string) *metrics.Gauge {
	return diagnostic.Gauge(name, description)
}

func MetricsToPrometheusWriter(writer io.Writer) {
	diagnostic.MetricsToWriter("prometheus", writer)
}
