// Package diagnostic implementation of error handler
package diagnostic

import (
	"io"
	"time"

	"github.com/sentinelos/tasker/internal/diagnostic/logger"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

func NewDiagnostic(o DiagnosticOptions) *Diagnostic {
	return &Diagnostic{DiagnosticOptions: o}
}

func (d *Diagnostic) Trace(message string, labels ...logger.Label) {
	d.writeLog(logger.Trace, message, labels...)
}

func (d *Diagnostic) Debug(message string, labels ...logger.Label) {
	d.writeLog(logger.Debug, message, labels...)
}

func (d *Diagnostic) Info(message string, labels ...logger.Label) {
	d.writeLog(logger.Info, message, labels...)
}

func (d *Diagnostic) Warn(message string, labels ...logger.Label) {
	d.writeLog(logger.Warn, message, labels...)
}

func (d *Diagnostic) Error(message string, labels ...logger.Label) {
	d.writeLog(logger.Error, message, labels...)
}

func (d *Diagnostic) Fatal(message string, labels ...logger.Label) {
	d.writeLog(logger.Fatal, message, labels...)
}

func (d *Diagnostic) Counter(name, description string) *metrics.Counter {
	return d.MetricsSet.Counter(name, description)
}

func (d *Diagnostic) Gauge(name, description string) *metrics.Gauge {
	return d.MetricsSet.Gauge(name, description)
}

func (d *Diagnostic) WritePrometheus(writer io.Writer) {
	d.MetricsSet.WritePrometheus(writer)
}

func (d *Diagnostic) writeLog(severity logger.Severity, message string, labels ...logger.Label) {
	if severity < d.Severity || len(d.LogWriters) == 0 {
		return
	}

	for _, writer := range d.LogWriters {
		writer.Write(d.Name, &logger.Entry{
			Severity: severity,
			Message:  message,
			Labels:   labels,
			Time:     time.Now(),
		})
	}
}
