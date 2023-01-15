package prometheus

// Prometheus defines a Prometheus writer.
type Prometheus struct {
	Options
}

type Options struct {
	Tags map[string]string
}
