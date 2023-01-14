package metrics

import (
	"io"
	"sync"
)

// Counter is a counter.
type Counter struct {
	v uint64
}

// Gauge is a gauge.
type Gauge struct {
	v uint64
}

type Metric interface {
	Write(writer io.Writer)
}

type NamedMetric struct {
	Name        string
	Description string
	Type        Type
	Labels      map[string]Metric
	isAux       bool
}

type Set struct {
	mu      sync.Mutex
	Metrics []*NamedMetric
	//summaries []*Summary
}
