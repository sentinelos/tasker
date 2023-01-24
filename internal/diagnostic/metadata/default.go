package metadata

var (
	_default = NewSet()
)

func Add(metadata *Metadata) {
	_default.Add(metadata)
}

func Get(name string) *Metadata {
	return _default.Get(name)
}
