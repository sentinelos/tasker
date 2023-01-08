// Package diagnostic implementation of error handler
package diagnostic

import (
	"time"

	"github.com/sentinelos/tasker/internal/diagnostic/logging"
)

func NewDiagnostic(o DiagnosticOptions) *Diagnostic {
	return &Diagnostic{DiagnosticOptions: o}
}

func (d *Diagnostic) Trace(message string, labels ...logging.Label) {
	d.writeLog(logging.Trace, message, labels...)
}

func (d *Diagnostic) Debug(message string, labels ...logging.Label) {
	d.writeLog(logging.Debug, message, labels...)
}

func (d *Diagnostic) Info(message string, labels ...logging.Label) {
	d.writeLog(logging.Info, message, labels...)
}

func (d *Diagnostic) Warn(message string, labels ...logging.Label) {
	d.writeLog(logging.Warn, message, labels...)
}

func (d *Diagnostic) Error(message string, labels ...logging.Label) {
	d.writeLog(logging.Error, message, labels...)
}

func (d *Diagnostic) Fatal(message string, labels ...logging.Label) {
	d.writeLog(logging.Fatal, message, labels...)
}

func (d *Diagnostic) Counter(name string, labels ...logging.Label) {

}

func (d *Diagnostic) writeLog(severity logging.Severity, message string, labels ...logging.Label) {
	if severity < d.Severity || len(d.LogWriters) == 0 {
		return
	}

	for _, writer := range d.LogWriters {
		writer.Write(d.Name, &logging.Entry{
			Severity: severity,
			Message:  message,
			Labels:   labels,
			Time:     time.Now(),
		})
	}
}
