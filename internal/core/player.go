package core

import (
	"context"
	"errors"
	"log"
	"sync"
)

// PLayer=orchestration only (SRP pattern).
type Player struct {
	reg Registry
	out Output
	pl  Playlist

	mu        sync.Mutex
	curIdx    int
	curTrack  *Track
	curStream DecodedStream
}

func NewPlayer(reg Registry, out Output, pl Playlist) *Player {
	return &Player{
		reg:    reg,
		out:    out,
		pl:     pl,
		curIdx: -1,
	}
}

func (p *Player) Play(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.curTrack != nil && p.out.State() == Playing {
		p.out.Resume()
		return nil
	}

	track, idx, ok := p.pl.Current()
	if !ok {
		return errors.New("playlist empty")
	}
	return p.playIndex(ctx, idx, track)
}

func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.out.Pause()
}

func (p *Player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.out.State() == Paused {
		p.out.Resume()
	}
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stopCurrent_locked()
}

func (p *Player) Next() {
	p.mu.Lock()
	defer p.mu.Unlock()

	track, idx, ok := p.pl.Next()
	if !ok {
		log.Println("End of playlist")
		return
	}

	_ = p.playIndex(context.Background(), idx, track)

}

func (p *Player) Previous() {
	p.mu.Lock()
	defer p.mu.Unlock()

	track, idx, ok := p.pl.Previous()
	if !ok {
		log.Println("Start of playlist")
		return
	}
	_ = p.playIndex(context.Background(), idx, track)
}

func (p *Player) playIndex(ctx context.Context, idx int, track Track) error {
	// Stop anything currently playing
	p.stopCurrent_locked()

	// Find decoder (LSP: any Decoder works here)
	dec := p.reg.FindDecoder(track.Path)
	if dec == nil {
		return errors.New("no decoder for " + track.Path)
	}

	stream, meta, err := dec.Decode(ctx, track.Path)
	if err != nil {
		return err
	}

	p.curIdx = idx
	p.curTrack = meta
	p.curStream = stream

	// Start playback; Output is responsible for async completion.
	p.out.Play(stream, func() {
		// auto-next when finished
		p.Next()
	})
	return nil
}

func (p *Player) stopCurrent_locked() {
	if p.curStream != nil {
		p.out.Stop()
		_ = p.curStream.Close()
		p.curStream = nil
	}
	p.curTrack = nil
	p.curIdx = -1
}
