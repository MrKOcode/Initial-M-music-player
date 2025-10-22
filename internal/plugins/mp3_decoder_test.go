package plugins

import (
	"context"
	"testing"
)

// to verify that mp3Decoder only accepts .mp3 files
func TestMP3Decoder_CanHandle(t *testing.T) {
	decoder := &mp3Decoder{}

	tests := []struct {
		path string
		want bool
	}{
		{"song.mp3", true},
		{"SONG.MP3", true},
		{"track.wav", false},
		{"audio.flac", false},
		{"noextension", false},
	}

	for _, tt := range tests {
		got := decoder.CanHandle(tt.path)
		if got != tt.want {
			t.Errorf("CanHandle(%q) = %v; want %v", tt.path, got, tt.want)
		}
	}
}

// to verify that Decode returns correct Track information
func TestMP3Decoder_Decode(t *testing.T) {
	decoder := &mp3Decoder{}
	ctx := context.Background()
	stream, track, err := decoder.Decode(ctx, "test_song.mp3")

	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}
	if stream == nil {
		t.Fatal("Decode returned nil stream")
	}
	if track == nil {
		t.Fatal("Decode returned nil track")
	}
	if track.Title != "test_song.mp3" {
		t.Errorf("Decode returned Title = %q; want %q", track.Title, "test_song.mp3")
	}
	if track.DurationSeconds != 5 {
		t.Errorf("Decode returned DurationSeconds = %d; want %d", track.DurationSeconds, 5)
	}

	// Ensure stream.Close works without error
	if err := stream.Close(); err != nil {
		t.Errorf("stream.Close() returned error: %v", err)
	}
}
