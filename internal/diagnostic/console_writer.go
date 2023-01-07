package diagnostic

import (
	"bytes"
	"os"
)

func NewConsoleWriter(o ConsoleWriterOptions) *ConsoleWriter {
	return &ConsoleWriter{ConsoleWriterOptions: o}
}

func (c *ConsoleWriter) WriteEntry(entry *Entry) {
	var buf bytes.Buffer

	buf.WriteString(entry.Time.Format(c.TimeFormat))
	buf.WriteString(" ")
	buf.WriteString(entry.Severity.String())

	if !c.EndWithMessage {
		buf.WriteString(" ")
		buf.WriteString(entry.Message)
	}

	for _, field := range entry.Fields {
		buf.WriteString(" ")
		buf.WriteString(field.Name)
		buf.WriteString("=")
		buf.WriteString(field.Value)
	}

	if c.EndWithMessage {
		buf.WriteString(" ")
		buf.WriteString(entry.Message)
	}

	buf.WriteString("\n")

	writer := os.Stdout
	if entry.Severity >= DiagError {
		writer = os.Stderr
	}

	buf.WriteTo(writer)
}
