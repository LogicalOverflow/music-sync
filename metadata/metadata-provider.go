package metadata

import (
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/util"
	"github.com/dhowden/tag"
	"os"
	"path/filepath"
)

// Provider is used to get metadata for songs
type Provider interface {
	CollectMetadata(song string) SongMetadata
}

// SongMetadata holds the metadata for a song
type SongMetadata struct {
	Title  string
	Artist string
	Album  string
}

func GetProvider() Provider {
	return basicProvider{}
}

type basicProvider struct{}

func (basicProvider) CollectMetadata(song string) SongMetadata {
	path := filepath.Join(playback.AudioDir, song)
	if !util.IsFile(path) {
		return SongMetadata{}
	}
	f, err := os.Open(path)
	if err != nil {
		return SongMetadata{}
	}
	md, err := tag.ReadFrom(f)
	if err != nil {
		return SongMetadata{}
	}
	return SongMetadata{
		Title:  md.Title(),
		Artist: md.Artist(),
		Album:  md.Album(),
	}
}
