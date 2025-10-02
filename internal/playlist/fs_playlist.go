package playlist

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"initial-m/internal/core"
)

// SRP: only concerns with listing and indexing tracks on filesystem.
type fsPlaylist struct {
	mu     sync.Mutex
	paths  []string
	cursor int
}

func NewFSPlaylist(folder string, exts []string) (core.Playlist, error) {
	allowed := map[string]bool{}
	for _, e := range exts {
		allowed[strings.ToLower(e)] = true
	}

	var files []string
	err := filepath.WalkDir(folder, func(p string, d os.DirEntry, err error) error {
		if err != nil { return err }
		if d.IsDir() { return nil }
		if allowed[strings.ToLower(filepath.Ext(p))] {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	pl := &fsPlaylist{paths: files, cursor: 0}
	return pl, nil
}

func (p *fsPlaylist) Len() int {
	p.mu.Lock(); defer p.mu.Unlock()
	return len(p.paths)
}

func (p *fsPlaylist) Current() (core.Track, int, bool) {
	p.mu.Lock(); defer p.mu.Unlock()
	if len(p.paths) == 0 { return core.Track{}, -1, false }
	return core.Track{Path: p.paths[p.cursor], Title: filepath.Base(p.paths[p.cursor]), DurationSeconds: 5}, p.cursor, true
}

func (p *fsPlaylist) Next() (core.Track, int, bool) {
	p.mu.Lock(); defer p.mu.Unlock()
	if len(p.paths) == 0 { return core.Track{}, -1, false }
	if p.cursor+1 >= len(p.paths) { return core.Track{}, -1, false }
	p.cursor++
	return core.Track{Path: p.paths[p.cursor], Title: filepath.Base(p.paths[p.cursor]), DurationSeconds: 5}, p.cursor, true
}

func (p *fsPlaylist) Previous() (core.Track, int, bool) {
	p.mu.Lock(); defer p.mu.Unlock()
	if len(p.paths) == 0 { return core.Track{}, -1, false }
	if p.cursor-1 < 0 { return core.Track{}, -1, false }
	p.cursor--
	return core.Track{Path: p.paths[p.cursor], Title: filepath.Base(p.paths[p.cursor]), DurationSeconds: 5}, p.cursor, true
}

func (p *fsPlaylist) JumpTo(i int) (core.Track, int, bool) {
	p.mu.Lock(); defer p.mu.Unlock()
	if i < 0 || i >= len(p.paths) { return core.Track{}, -1, false }
	p.cursor = i
	return core.Track{Path: p.paths[p.cursor], Title: filepath.Base(p.paths[p.cursor]), DurationSeconds: 5}, p.cursor, true
}
