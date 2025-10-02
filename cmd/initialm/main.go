package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"initial-m/internal/core"
	"initial-m/internal/output"
	"initial-m/internal/playlist"
	"initial-m/internal/plugins"
)

func main() {
	// 1) Wire dependencies (DIP): everything by interfaces.
	reg := core.NewRegistry()
	plugins.RegisterMP3Decoder(reg) // OCP: add more by registering.

	out := output.NewSimOutput() // swap with a real output later

	// 2) Build playlist from a folder given by arg or current dir
	folder := "."
	if len(os.Args) > 1 {
		folder = os.Args[1]
	}
	pl, err := playlist.NewFSPlaylist(folder, []string{".mp3"})
	if err != nil {
		log.Fatal(err)
	}
	if pl.Len() == 0 {
		log.Printf("No .mp3 files found in %s", folder)
	}

	// 3) Create player
	p := core.NewPlayer(reg, out, pl)

	// 4) Control transport (simple demo loop)
	go func() {
		time.Sleep(500 * time.Millisecond)
		p.Play(context.Background())

		time.Sleep(2 * time.Second)
		p.Pause()

		time.Sleep(1 * time.Second)
		p.Resume()

		time.Sleep(2 * time.Second)
		p.Next()

		time.Sleep(2 * time.Second)
		p.Previous()

		time.Sleep(3 * time.Second)
		p.Stop()
	}()

	// 5) graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("Exiting Initial M...")
	_ = filepath.Base("") // keep import happy if unused on some Go versions
}

