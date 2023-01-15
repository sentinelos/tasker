package metrics

import (
	"fmt"
	"strings"
)

// NewSet creates new set of metrics.
//
// Pass the set to RegisterSet() function in order to export its metrics via global WritePrometheus() call.
func NewSet() *Set {
	return &Set{
		Metrics: []*NamedMetric{},
	}
}

// Counter registers and returns new counter with the given name and description in the s.
//
// name must be valid Prometheus-compatible Metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
//
// The returned counter is safe to use from concurrent goroutines.
func (s *Set) Counter(name, description string) *Counter {
	return s.Metric(name, description, &Counter{}, false).(*Counter)
}

// Gauge registers and returns gauge with the given name and description in s, which calls f
// to obtain gauge value.
//
// name must be valid Prometheus-compatible Metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
//
// The returned gauge is safe to use from concurrent goroutines.
func (s *Set) Gauge(name, description string) *Gauge {
	return s.Metric(name, description, &Gauge{}, false).(*Gauge)
}

// Metric registers given Metric with the given name and description.
func (s *Set) Metric(name, description string, metric Metric, isAux bool) Metric {
	if err := validateMetric(name); err != nil {
		panic(fmt.Errorf("invalid metric name %q: %s", name, err))
	}

	index := strings.IndexByte(name, '{')
	n := name
	labels := ""

	if index > 0 {
		n = name[:index]
		labels = name[index+1 : len(name)-1]
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, m := range s.Metrics {
		if met, found := m.Labels[labels]; m.Name == n && found {
			return met
		}
	}

	for _, m := range s.Metrics {
		if m.Name == n {
			m.Labels[labels] = metric
			return metric
		}
	}

	var t Type
	switch metric.(type) {
	case *Counter:
		t = TypeCounter
	case *Gauge:
		t = TypeGauge
	}

	s.Metrics = append(s.Metrics, &NamedMetric{
		Name:        n,
		Description: description,
		Type:        t,
		Labels: map[string]Metric{
			labels: metric,
		},
		isAux: isAux,
	})

	return metric
}
