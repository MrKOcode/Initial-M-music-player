package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"initial-m/internal/core"
	"initial-m/internal/output"
	"initial-m/internal/playlist"
	"initial-m/internal/plugins"
)

type Server struct {
	player *core.Player
}

func (s *Server) handlePlay(w http.ResponseWriter, r *http.Request) {
	err := s.player.Play(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "playing"})
}

func (s *Server) handlePause(w http.ResponseWriter, r *http.Request) {
	s.player.Pause()
	json.NewEncoder(w).Encode(map[string]string{"status": "paused"})
}

func (s *Server) handleResume(w http.ResponseWriter, r *http.Request) {
	s.player.Resume()
	json.NewEncoder(w).Encode(map[string]string{"status": "resumed"})
}

func (s *Server) handleNext(w http.ResponseWriter, r *http.Request) {
	s.player.Next()
	json.NewEncoder(w).Encode(map[string]string{"status": "next"})
}

func (s *Server) handlePrev(w http.ResponseWriter, r *http.Request) {
	s.player.Previous()
	json.NewEncoder(w).Encode(map[string]string{"status": "previous"})
}

func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	s.player.Stop()
	json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
}

func main() {
	// Init registry, output, playlist
	reg := core.NewRegistry()
	plugins.RegisterMP3Decoder(reg)
	out := output.NewSimOutput()

	folder := "."
	if len(os.Args) > 1 {
		folder = os.Args[1]
	}
	pl, err := playlist.NewFSPlaylist(folder, []string{".mp3"})
	if err != nil {
		log.Fatal(err)
	}

	player := core.NewPlayer(reg, out, pl)
	server := &Server{player: player}

	// REST routes
	http.HandleFunc("/play", server.handlePlay)
	http.HandleFunc("/pause", server.handlePause)
	http.HandleFunc("/resume", server.handleResume)
	http.HandleFunc("/next", server.handleNext)
	http.HandleFunc("/prev", server.handlePrev)
	http.HandleFunc("/stop", server.handleStop)

	log.Println("🎵 Initial M REST API running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
