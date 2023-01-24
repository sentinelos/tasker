package diagnostic

import (
	"io"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

var (
	_default = &Diagnostic{}
)

func init() {
	_default = NewDiagnostic(NewOptions(
		WithName("default"),
	))
}

func Trace(message string) {
	_default.Trace(message, labels.Labels{})
}

func TraceWithLabels(message string, labels labels.Labels) {
	_default.Trace(message, labels)
}

func Debug(message string) {
	_default.Debug(message, labels.Labels{})
}

func DebugWithLabels(message string, labels labels.Labels) {
	_default.Debug(message, labels)
}

func Info(message string) {
	_default.Info(message, labels.Labels{})
}

func InfoWithLabels(message string, labels labels.Labels) {
	_default.Info(message, labels)
}

func Warn(message string) {
	_default.Warn(message, labels.Labels{})
}

func WarnWithLabels(message string, labels labels.Labels) {
	_default.Warn(message, labels)
}

func Error(message string) {
	_default.Error(message, labels.Labels{})
}

func ErrorWithLabels(message string, labels labels.Labels) {
	_default.Error(message, labels)
}

func Fatal(message string) {
	_default.Fatal(message, labels.Labels{})
}

func FatalWithLabels(message string, labels labels.Labels) {
	_default.Fatal(message, labels)
}

func Counter(name, description string) *metrics.CounterLabels {
	return _default.Counter(name, description)
}

func Gauge(name, description string) *metrics.GaugeLabels {
	return _default.Gauge(name, description)
}

func MetricsToPrometheusWriter(writer io.Writer) {
	_default.MetricsToWriter("prometheus", writer)
}
