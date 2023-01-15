// Package prometheus implementation of diagnostic metrics writer to prometheus
package prometheus

import (
	"fmt"
	"io"
	"strings"

	"github.com/sentinelos/tasker/internal/diagnostic/metrics"
)

func NewPrometheus(o Options) *Prometheus {
	return &Prometheus{Options: o}
}

func (p *Prometheus) Write(writer io.Writer, set *metrics.Set) {
	var tags []string
	for tag, value := range p.Tags {
		tags = append(tags, tag+"=\""+value+"\"")
	}

	for _, metric := range set.Metrics {
		fmt.Fprintf(writer, "# HELP %s %s\n", metric.Name, metric.Description)
		fmt.Fprintf(writer, "# TYPE %s %s\n", metric.Name, metric.Type.String())

		for label, met := range metric.Labels {
			labels := tags
			if len(label) > 0 {
				labels = append(labels, label)
			}

			fmt.Fprintf(writer, "%s{%s} ", metric.Name, strings.Join(labels, ", "))
			met.Write(writer)
			fmt.Fprintf(writer, "\n")
		}
	}
}
