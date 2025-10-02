package core

// registry implements the Registry interface.
type registry struct {
	decoders []Decoder
}

func NewRegistry() Registry {
	return &registry{}
}

func (r *registry) RegisterDecoder(d Decoder) {
	r.decoders = append(r.decoders, d)
}

func (r *registry) FindDecoder(path string) Decoder {
	for _, d := range r.decoders {
		if d.CanHandle(path) {
			return d
		}
	}
	return nil
}
