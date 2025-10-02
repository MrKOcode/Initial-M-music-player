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


//DecodedStream is intentionally opaque to the player
//The output knows how to handle it
type DecodedStream interface {
	Close() error
}

//output plays a decoded stream
type Output interface {
	Play(stream DecodedStream, onEnd func()) 
	Pause()
	Resume()
	Stop()
	State() OutputState
}

type OutputState int

const (
	Stopped OutputState = iota
	Playing
	Paused
)

//playlist abstracts a list of tracks
type Playlist interface {
	Current() (Track, int, bool)
	Next() (Track, int, bool)
	Previous() (Track, int, bool)
	JumpTo(i int) (Track, int, bool)
	Len() int
}

//Registry for playable plugins
type Registry interface {
	RegisterDecoder(Decoder)	
	FindDecoder(path string) Decoder
}



