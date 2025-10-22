// internal/core/registry_test.go
package core_test

import (
    "testing"

    "initial-m/internal/core"
    "initial-m/internal/plugins"
)

func TestRegistry_FindDecoder(t *testing.T) {
    reg := core.NewRegistry()

    // Register MP3 adapter
    plugins.RegisterMP3Decoder(reg)

    tests := []struct {
        name      string
        path      string
        wantFound bool
    }{
        {"mp3 file", "song.mp3", true},
        {"uppercase MP3", "SONG.MP3", true},
        {"flac file", "album.flac", false},
        {"wav file", "drums.wav", false},
        {"no extension", "track", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            d := reg.FindDecoder(tt.path)
            if (d != nil) != tt.wantFound {
                t.Errorf("FindDecoder(%q) found = %v; want %v", tt.path, d != nil, tt.wantFound)
            }
        })
    }
}
