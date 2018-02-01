package metadata

import (
	"encoding/json"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/util"
	"os"
	"path/filepath"
)

// LyricsAtom describes a single word/element in the lyrics
type LyricsAtom struct {
	Timestamp int64  `json:"timestamp"`
	Caption   string `json:"caption"`
}

// LyricsLine is a line in the lyrics, composed of LyricsAtoms
type LyricsLine []LyricsAtom

// LyricsProvider is used to get the lyrics for songs
type LyricsProvider interface {
	CollectLyrics(song string) []LyricsLine
}

// GetLyricsProvider returns a new LyricsProvider
func GetLyricsProvider() LyricsProvider {
	return basicLyricsProvider{}
}

type basicLyricsProvider struct{}

func (basicLyricsProvider) CollectLyrics(song string) []LyricsLine {
	path := filepath.Join(playback.AudioDir, song+".json")
	if !util.IsFile(path) {
		return []LyricsLine{}
	}
	f, err := os.Open(path)
	if err != nil {
		return []LyricsLine{}
	}
	result := make([]LyricsLine, 0)

	if err := json.NewDecoder(f).Decode(&result); err != nil {
		return []LyricsLine{}
	}

	return result
}
