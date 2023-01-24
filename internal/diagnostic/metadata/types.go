package metadata

import (
	"sync"

	"github.com/sentinelos/tasker/internal/diagnostic/labels"
)

type (
	// Detector define contract for the metadata detection.
	Detector func() labels.Labels

	Metadata struct {
		Name        string
		Description string
		Labels      labels.Labels
	}

	Set struct {
		mu       sync.Mutex
		Metadata []*Metadata
	}
)
