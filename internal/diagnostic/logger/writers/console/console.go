// Package console implementation of diagnostic logger writer to console
package console

import (
	"bytes"
	"os"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
	"github.com/sentinelos/tasker/internal/diagnostic/logger"
)

func NewConsole(o Options) *Console {
	return &Console{Options: o}
}

func (c *Console) Write(entry *logger.Entry) {
	var buf bytes.Buffer

	buf.WriteString("[")
	buf.WriteString(entry.Stamp.Format(c.TimeFormat))
	buf.WriteString(" - ")
	buf.WriteString(entry.Severity.String())
	buf.WriteString("]")

	if !c.EndWithMessage {
		buf.WriteString(" \"")
		buf.WriteString(entry.Message)
		buf.WriteString("\"")
	}

	buf.WriteString("{")
	buf.WriteString(labels.Merge(c.Tags, entry.Labels).String())
	buf.WriteString("}")

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
