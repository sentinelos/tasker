package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Increment increment Counter value.
func (c *Counter) Increment() {
	c.Add(1)
}

// Reset sets Counter value to 0.
func (c *Counter) Reset() {
	c.Set(0)
}

// Add sets Counter value to v.
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

// Write writes Counter value to writers.
func (c *Counter) Write(writer io.Writer) {
	fmt.Fprintf(writer, "%d", c.Get())
}
