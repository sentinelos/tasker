// Package diagnostic implementation of error handler
package diagnostic

import (
	"fmt"
	"io"
	"time"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

func NewDiagnostic(o Options) *Diagnostic {
	return &Diagnostic{Options: o}
}

func (d *Diagnostic) Trace(message string, labels labels.Labels) {
	d.LoggerToWriter(logger.Trace, message, labels)
}

func (d *Diagnostic) Debug(message string, labels labels.Labels) {
	d.LoggerToWriter(logger.Debug, message, labels)
}

func (d *Diagnostic) Info(message string, labels labels.Labels) {
	d.LoggerToWriter(logger.Info, message, labels)
}

func (d *Diagnostic) Warn(message string, labels labels.Labels) {
	d.LoggerToWriter(logger.Warn, message, labels)
}

func (d *Diagnostic) Error(message string, labels labels.Labels) {
	d.LoggerToWriter(logger.Error, message, labels)
}

func (d *Diagnostic) Fatal(message string, labels labels.Labels) {
	d.LoggerToWriter(logger.Fatal, message, labels)
}

func (d *Diagnostic) Counter(name, description string) *metrics.CounterLabels {
	return d.MetricsSet.Counter(name, description)
}

func (d *Diagnostic) Gauge(name, description string) *metrics.GaugeLabels {
	return d.MetricsSet.Gauge(name, description)
}

func (d *Diagnostic) LoggerToWriter(severity logger.Severity, message string, labels labels.Labels) {
	d.loggerToMetric(severity, labels)

	if severity < d.Severity || len(d.LoggerWriters) == 0 {
		return
	}

	stamp := time.Now()
	for _, writer := range d.LoggerWriters {
		writer.Write(&logger.Entry{
			Severity: severity,
			Message:  message,
			Labels:   labels,
			Stamp:    stamp,
		})
	}
}

func (d *Diagnostic) loggerToMetric(severity logger.Severity, labels labels.Labels) {
	name := fmt.Sprintf("diagnostic_%s_total", severity.String())
	description := fmt.Sprintf("Total number of diagnostic logger %ss", severity.String())

	d.Counter(name, description).Get(labels).Increment()
}

func (d *Diagnostic) MetricsToWriter(name string, writer io.Writer) {
	if metricsWriter, found := d.MetricsWriters[name]; found {
		metricsWriter.Write(writer, d.MetricsSet)
	}
}
