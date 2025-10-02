package plugins

import (
	"context"
	"path/filepath"
	"strings"

	"initial-m/internal/core"
)

// LSP/OCP: This plugin "looks" like a decoder; Player never needs to know which one.
// This demo just returns a dummy stream to keep things runnable without 3rd-party libs.

type mp3Decoder struct{}

func (m *mp3Decoder) CanHandle(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".mp3")
}

func (m *mp3Decoder) Decode(ctx context.Context, path string) (core.DecodedStream, *core.Track, error) {
	// In a real implementation, open file, decode frames, return an object the Output can pull from.
	// Here we just simulate based on filename length.
	dur := 5
	title := filepath.Base(path)
	stream := &dummyStream{path: path}
	return stream, &core.Track{Path: path, Title: title, DurationSeconds: dur}, nil
}

type dummyStream struct{ path string }

func (d *dummyStream) Close() error { return nil }

// Registry wiring
type registry struct{ decoders []core.Decoder }

func NewRegistry() core.Registry { return &registry{} }

func (r *registry) RegisterDecoder(d core.Decoder) { r.decoders = append(r.decoders, d) }

func (r *registry) FindDecoder(path string) core.Decoder {
	for _, d := range r.decoders {
		if d.CanHandle(path) {
			return d
		}
	}
	return nil
}

func RegisterMP3Decoder(reg core.Registry) {
	reg.RegisterDecoder(&mp3Decoder{})
}
