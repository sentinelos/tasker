package diagnostic

import (
	"io"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
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

func Trace(message string) {
	diagnostic.Trace(message, labels.Labels{})
}

func TraceWithLabels(message string, labels labels.Labels) {
	diagnostic.Trace(message, labels)
}

func Debug(message string) {
	diagnostic.Debug(message, labels.Labels{})
}

func DebugWithLabels(message string, labels labels.Labels) {
	diagnostic.Debug(message, labels)
}

func Info(message string) {
	diagnostic.Info(message, labels.Labels{})
}

func InfoWithLabels(message string, labels labels.Labels) {
	diagnostic.Info(message, labels)
}

func Warn(message string) {
	diagnostic.Warn(message, labels.Labels{})
}

func WarnWithLabels(message string, labels labels.Labels) {
	diagnostic.Warn(message, labels)
}

func Error(message string) {
	diagnostic.Error(message, labels.Labels{})
}

func ErrorWithLabels(message string, labels labels.Labels) {
	diagnostic.Error(message, labels)
}

func Fatal(message string) {
	diagnostic.Fatal(message, labels.Labels{})
}

func FatalWithLabels(message string, labels labels.Labels) {
	diagnostic.Fatal(message, labels)
}

func Counter(name, description string) *metrics.CounterLabels {
	return diagnostic.Counter(name, description)
}

func Gauge(name, description string) *metrics.GaugeLabels {
	return diagnostic.Gauge(name, description)
}

func MetricsToPrometheusWriter(writer io.Writer) {
	diagnostic.MetricsToWriter("prometheus", writer)
}
