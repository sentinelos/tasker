package metrics

// Type represents the metric type.
type Type uint

const (
	TypeCounter Type = iota
	TypeGauge
)

var (
	TypeNames = map[Type]string{
		TypeCounter: "counter",
		TypeGauge:   "gauge",
	}
)

// String severity to string
func (t Type) String() string {
	return TypeNames[t]
}
