package metrics

import (
	"fmt"
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
//   - foo_bar
//
// The returned counter is safe to use from concurrent goroutines.
func (s *Set) Counter(name, description string) *CounterLabels {
	return s.getOrCreate(name, description, TypeCounter, false).(*CounterLabels)
}

// Gauge registers and returns gauge with the given name and description in s, which calls f
// to obtain gauge value.
//
// name must be valid Prometheus-compatible Metric with possible labels.
// For instance,
//
//   - foo
//   - foo_bar
//
// The returned gauge is safe to use from concurrent goroutines.
func (s *Set) Gauge(name, description string) *GaugeLabels {
	return s.getOrCreate(name, description, TypeGauge, false).(*GaugeLabels)
}

func (s *Set) getOrCreate(name, description string, mType Type, isAux bool) Labels {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validateMetric(name); err != nil {
		panic(fmt.Errorf("invalid metric name %q: %s", name, err))
	}

	if index := s.Index(name, mType); index > -1 {
		return s.Metrics[index].Labels
	}

	var labels Labels

	switch mType {
	case TypeCounter:
		labels = &CounterLabels{}
	case TypeGauge:
		labels = &GaugeLabels{}
	}

	s.Metrics = append(s.Metrics, &NamedMetric{
		Name:        name,
		Description: description,
		Type:        mType,
		Labels:      labels,
		isAux:       isAux,
	})

	return labels
}

func (s *Set) Index(name string, mType Type) int {
	for index, m := range s.Metrics {
		if m.Name == name && m.Type == mType {
			return index
		}
	}

	return -1
}
