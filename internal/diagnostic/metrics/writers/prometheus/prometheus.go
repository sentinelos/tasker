// Package prometheus implementation of diagnostic metrics writer to prometheus
package prometheus

import (
	"fmt"
	"io"

	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

func NewPrometheus(o Options) *Prometheus {
	return &Prometheus{Options: o}
}

func (p *Prometheus) Write(writer io.Writer, set *metrics.Set) {
	tags := p.Tags.String()

	for _, metric := range set.Metrics {
		fmt.Fprintf(writer, "# HELP %s %s\n", metric.Name, metric.Description)
		fmt.Fprintf(writer, "# TYPE %s %s\n", metric.Name, metric.Type.String())

		for label, met := range metric.Labels.Values() {
			if len(label) > 0 {
				fmt.Fprintf(writer, "%s{%s,%s} ", metric.Name, tags, label)
			} else {
				fmt.Fprintf(writer, "%s{%s} ", metric.Name, tags)
			}
			met.Write(writer)
			fmt.Fprintf(writer, "\n")
		}
	}
}
