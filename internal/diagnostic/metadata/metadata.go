// Package metadata implementation of diagnostic metadata
package metadata

func NewMetadata(name, description string, detector Detector) *Metadata {
	return &Metadata{
		Name:        name,
		Description: description,
		Labels:      detector(),
	}
}

// GetLabel returns the value in the map for the provided label
func (md *Metadata) GetLabel(name string) string {
	return md.Labels.Get(name)
}
