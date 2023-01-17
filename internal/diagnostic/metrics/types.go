package metrics

import (
	"io"
	"sync"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

type Metric interface {
	Write(writer io.Writer)
}

// Labels allows you to present labels independently of their storage.
type Labels interface {
	// Has returns whether the provided label exists.
	Has(labels labels.Labels) bool

	// Delete deletes the Metric label for the provided label.
	Delete(labels labels.Labels)

	// Values returns the values of Metric labels.
	Values() map[string]Metric
}

// Counter is a counter metric
type Counter struct {
	v uint64
}

// Gauge is a gauge metric
type Gauge struct {
	v uint64
}

type NamedMetric struct {
	Name        string
	Description string
	Type        Type
	Labels      Labels
	isAux       bool
}

type Set struct {
	mu      sync.Mutex
	Metrics []*NamedMetric
}
