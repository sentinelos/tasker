package metadata

// NewSet creates new set of metadata.
func NewSet() *Set {
	return &Set{
		Metadata: []*Metadata{},
	}
}

func (s *Set) Add(metadata *Metadata) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index := s.Index(metadata.Name); index > -1 {
		return
	}

	s.Metadata = append(s.Metadata, metadata)
}

func (s *Set) Get(name string) *Metadata {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index := s.Index(name); index > -1 {
		return s.Metadata[index]
	}

	return &Metadata{}
}

func (s *Set) Index(name string) int {
	for index, m := range s.Metadata {
		if m.Name == name {
			return index
		}
	}

	return -1
}
