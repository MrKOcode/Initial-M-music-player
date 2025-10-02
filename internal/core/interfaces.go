package core

import "context"

//track drscribes a playable item
type Track struct {
	Title           string
	Path            string
	DurationSeconds int
}

//decoder decodes a given file into tracks
type Decoder interface {
	CanHandle(path string) bool
	//Decode produces a decodedStream and basic matadata
	Decoder(ctx context.Context, path string) (DecodedStream, Track, error)
}
