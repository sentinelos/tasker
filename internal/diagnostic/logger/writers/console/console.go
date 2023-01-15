// Package console implementation of diagnostic logger writer to console
package console

import (
	"bytes"
	"os"

	"github.com/sentinelos/tasker/internal/diagnostic/logger"
)

func NewConsole(o Options) *Console {
	return &Console{Options: o}
}

func (c *Console) Write(entry *logger.Entry) {
	var buf bytes.Buffer

	buf.WriteString("[")
	buf.WriteString(entry.Time.Format(c.TimeFormat))
	buf.WriteString(" - ")
	buf.WriteString(entry.Severity.String())
	buf.WriteString("]")

	if !c.EndWithMessage {
		buf.WriteString(" \"")
		buf.WriteString(entry.Message)
		buf.WriteString("\"")
	}

	for tag, value := range c.Tags {
		buf.WriteString(" ")
		buf.WriteString(tag)
		buf.WriteString("=")
		buf.WriteString(value)
	}

	for _, field := range entry.Labels {
		buf.WriteString(" ")
		buf.WriteString(field.Name)
		buf.WriteString("=")
		buf.WriteString(field.Value)
	}

	if c.EndWithMessage {
		buf.WriteString(" \"")
		buf.WriteString(entry.Message)
		buf.WriteString("\"")
	}

	buf.WriteString("\n")

	writer := os.Stdout
	if entry.Severity >= logger.Error {
		writer = os.Stderr
	}

	buf.WriteTo(writer)
}
