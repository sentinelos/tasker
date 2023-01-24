package metrics

import (
	"fmt"
	"io"
	"sync/atomic"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

// Increment increment Gauge value.
func (g *Gauge) Increment() {
	g.Add(1)
}

// Decrement decrement Gauge value.
func (g *Gauge) Decrement() {
	g.Add(g.Get() - 1)
}

// Add adds Gauge value to v.
func (g *Gauge) Add(v uint64) {
	atomic.AddUint64(&g.v, v)
}

// Set sets Gauge value to v.
func (g *Gauge) Set(v uint64) {
	atomic.StoreUint64(&g.v, v)
}

// Get returns the current value for Gauge.
func (g *Gauge) Get() uint64 {
	return atomic.LoadUint64(&g.v)
}

// Write writes Gauge value to writer.
func (g *Gauge) Write(writer io.Writer) {
	fmt.Fprintf(writer, "%d", g.Get())
}

// Get returns the Gauge for the provided label.
func (gl GaugeLabels) Get(labels labels.Labels) *Gauge {
	return gl[labels.String()]
}

// Has returns whether the provided label exists.
func (gl GaugeLabels) Has(label labels.Labels) bool {
	_, exists := gl[label.String()]

	return exists
}

// Delete deletes the Gauge label for the provided label.
func (gl GaugeLabels) Delete(labels labels.Labels) {
	delete(gl, labels.String())
}

// Values returns the values of GaugeLabels.
func (gl GaugeLabels) Values() map[string]Metric {
	var values = map[string]Metric{}

	for label, value := range gl {
		values[label] = value
	}

	return values
}
