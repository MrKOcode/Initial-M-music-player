package output

import (
	"log"
	"sync"
	"time"

	"initial-m/internal/core"
)

// SRP: just "speaks" to an audio device (simulated).
// DIP: depends on core.DecodedStream, not on concrete plugin types.

type simOutput struct {
	mu    sync.Mutex
	state core.OutputState
	stop  chan struct{}
}

func NewSimOutput() core.Output { return &simOutput{state: core.Stopped} }

func (o *simOutput) Play(stream core.DecodedStream, onEnd func()) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// stop any prior
	if o.stop != nil {
		close(o.stop)
	}

	o.stop = make(chan struct{})
	o.state = core.Playing

	// Simulate playback length with a ticker. In real code you'd pull PCM frames.
	go func(done chan struct{}) {
		log.Println("Playing (sim):", stream)
		t := time.NewTimer(5 * time.Second) // pretend every track lasts ~5s
		select {
		case <-t.C:
			o.mu.Lock()
			o.state = core.Stopped
			o.mu.Unlock()
			if onEnd != nil {
				onEnd()
			}
		case <-done:
			t.Stop()
		}
	}(o.stop)
}

func (o *simOutput) Pause() {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.state == core.Playing {
		o.state = core.Paused
		log.Println("Paused")
	}
}

func (o *simOutput) Resume() {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.state == core.Paused {
		o.state = core.Playing
		log.Println("Resumed")
	}
}

func (o *simOutput) Stop() {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.stop != nil {
		close(o.stop)
		o.stop = nil
	}
	o.state = core.Stopped
	log.Println("Stopped")
}

func (o *simOutput) State() core.OutputState {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.state
}
