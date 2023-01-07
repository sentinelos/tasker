package diagnostic

import (
	"bytes"
	"io"
	"os"
)

type ConsoleWriter struct {
	// ColorOutput determines if used colorized output.
	ColorOutput bool

	// QuoteString determines if quoting string values.
	QuoteString bool

	// EndWithMessage determines if output message in the end.
	EndWithMessage bool

	// Writer is the output destination. using os.Stderr if empty.
	Writer io.Writer
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{
		ColorOutput:    true,
		QuoteString:    true,
		EndWithMessage: true,
		Writer:         os.Stderr,
	}
}

func (c *ConsoleWriter) WriteEntry(entry *Entry) {
	var buf bytes.Buffer

	buf.WriteString(entry.Time.Format(DefaultTimeStamp))
	buf.WriteString(" ")
	buf.WriteString(entry.Message)
	buf.WriteString(" ")
	buf.WriteString("severity")
	buf.WriteString("=")
	buf.WriteString(entry.Severity.String())
	for _, field := range entry.Fields {
		buf.WriteString(" ")
		buf.WriteString(field.Name)
		buf.WriteString("=")
		buf.WriteString(field.Value)
	}

	buf.WriteTo(c.Writer)
}
