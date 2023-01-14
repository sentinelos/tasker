// Package console implementation of diagnostic log writer to console
package console

import (
	"bytes"
	"os"

	"github.com/sentinelos/tasker/internal/diagnostic/logger"
)

func NewConsole(o ConsoleOptions) *Console {
	return &Console{ConsoleOptions: o}
}

func (c *Console) Write(name string, entry *logger.Entry) {
	var buf bytes.Buffer

	buf.WriteString("[")
	buf.WriteString(entry.Time.Format(c.TimeFormat))
	buf.WriteString("] [")
	buf.WriteString(name)
	buf.WriteString("] [")
	buf.WriteString(entry.Severity.String())
	buf.WriteString("]")

	if !c.EndWithMessage {
		buf.WriteString(" \"")
		buf.WriteString(entry.Message)
		buf.WriteString("\"")
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
