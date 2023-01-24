package metrics

import (
	"fmt"
	"io"
	"sync/atomic"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

// Increment increment Counter value.
func (c *Counter) Increment() {
	c.Add(1)
}

// Reset sets Counter value to 0.
func (c *Counter) Reset() {
	c.Set(0)
}

// Add adds Counter value to v.
func (c *Counter) Add(v uint64) {
	atomic.AddUint64(&c.v, v)
}

// Set sets Counter value to v.
func (c *Counter) Set(v uint64) {
	atomic.StoreUint64(&c.v, v)
}

// Get returns the current value for Counter.
func (c *Counter) Get() uint64 {
	return atomic.LoadUint64(&c.v)
}

// Write writes Counter value to writer.
func (c *Counter) Write(writer io.Writer) {
	fmt.Fprintf(writer, "%d", c.Get())
}

// Get returns the Counter for the provided label.
func (cl CounterLabels) Get(labels labels.Labels) *Counter {
	if c, found := cl[labels.String()]; found {
		return c
	}

	counter := &Counter{}
	cl[labels.String()] = counter

	return counter
}

// Has returns whether the provided label exists.
func (cl CounterLabels) Has(labels labels.Labels) bool {
	_, exists := cl[labels.String()]

	return exists
}

// Delete deletes the Counter label for the provided label.
func (cl CounterLabels) Delete(labels labels.Labels) {
	delete(cl, labels.String())
}

// Values returns the values of CounterLabels.
func (cl CounterLabels) Values() map[string]Metric {
	var values = map[string]Metric{}

	for label, value := range cl {
		values[label] = value
	}

	return values
}
