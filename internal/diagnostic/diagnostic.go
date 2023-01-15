// Package diagnostic implementation of error handler
package diagnostic

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sentinelos/tasker/internal/diagnostic/logger"
	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

func NewDiagnostic(o Options) *Diagnostic {
	return &Diagnostic{Options: o}
}

func (d *Diagnostic) Trace(message string, labels ...logger.Label) {
	d.ToLoggerWriter(logger.Trace, message, labels...)
}

func (d *Diagnostic) Debug(message string, labels ...logger.Label) {
	d.ToLoggerWriter(logger.Debug, message, labels...)
}

func (d *Diagnostic) Info(message string, labels ...logger.Label) {
	d.ToLoggerWriter(logger.Info, message, labels...)
}

func (d *Diagnostic) Warn(message string, labels ...logger.Label) {
	d.ToLoggerWriter(logger.Warn, message, labels...)
}

func (d *Diagnostic) Error(message string, labels ...logger.Label) {
	d.ToLoggerWriter(logger.Error, message, labels...)
}

func (d *Diagnostic) Fatal(message string, labels ...logger.Label) {
	d.ToLoggerWriter(logger.Fatal, message, labels...)
}

func (d *Diagnostic) Counter(name, description string) *metrics.Counter {
	return d.MetricsSet.Counter(name, description)
}

func (d *Diagnostic) Gauge(name, description string) *metrics.Gauge {
	return d.MetricsSet.Gauge(name, description)
}

func (d *Diagnostic) ToLoggerWriter(severity logger.Severity, message string, labels ...logger.Label) {
	d.loggerToMetric(severity, labels...)

	if severity < d.Severity || len(d.LoggerWriters) == 0 {
		return
	}

	for _, writer := range d.LoggerWriters {
		writer.Write(&logger.Entry{
			Severity: severity,
			Message:  message,
			Labels:   labels,
			Time:     time.Now(),
		})
	}
}

func (d *Diagnostic) loggerToMetric(severity logger.Severity, labels ...logger.Label) {
	var ls []string

	for _, label := range labels {
		ls = append(ls, label.Name+"=\""+label.Value+"\"")
	}

	metricName := fmt.Sprintf("diagnostic_%s_total{%s}", severity.String(), strings.Join(ls, ", "))
	metricDescription := fmt.Sprintf("Total number of diagnostic logger %ss", severity.String())

	d.Counter(metricName, metricDescription).Increment()
}

func (d *Diagnostic) MetricsToWriter(name string, writer io.Writer) {
	if metricsWriter, found := d.MetricsWriters[name]; found {
		metricsWriter.Write(writer, d.MetricsSet)
	}
}
