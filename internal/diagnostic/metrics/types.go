package metrics

import (
	"io"
	"sync"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

type (
	Metric interface {
		Write(writer io.Writer)
	}

	// Labels allows you to present labels independently of their storage.
	Labels interface {
		// Has returns whether the provided label exists.
		Has(labels labels.Labels) bool

		// Delete deletes the Metric label for the provided label.
		Delete(labels labels.Labels)

		// Values returns the values of Metric labels.
		Values() map[string]Metric
	}

	// Counter is a counter metric
	Counter struct {
		v uint64
	}

	CounterLabels map[string]*Counter

	// Gauge is a gauge metric
	Gauge struct {
		v uint64
	}

	GaugeLabels map[string]*Gauge

	// Type represents the Metric type.
	Type uint

	NamedMetric struct {
		Name        string
		Description string
		Type        Type
		Labels      Labels
		isAux       bool
	}

	Set struct {
		mu      sync.Mutex
		Metrics []*NamedMetric
	}
)
