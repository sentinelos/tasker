package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Increment increment Gauge value.
func (g *Gauge) Increment() {
	g.Add(1)
}

// Decrement decrement Gauge value.
func (g *Gauge) Decrement() {
	g.Add(g.Get() - 1)
}

// Add sets Gauge value to v.
func (g *Gauge) Add(v uint64) {
	atomic.AddUint64(&g.v, v)
}

// Set sets Gauge value to v.
func (g *Gauge) Set(v uint64) {
	atomic.StoreUint64(&g.v, v)
}

// Get returns the current value for c.
func (g *Gauge) Get() uint64 {
	return atomic.LoadUint64(&g.v)
}

// Write writes c to writer.
func (g *Gauge) Write(writer io.Writer) {
	fmt.Fprintf(writer, "%d", g.Get())
}
