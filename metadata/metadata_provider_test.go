package metadata

import (
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetProvider(t *testing.T) {
	assert.NotNil(t, GetProvider(), "GetProvider returned nil")
}

func TestBasicProvider_CollectMetadata(t *testing.T) {
	playback.AudioDir = "_test_files"
	bp := basicProvider{}

	assert.Equal(t, SongMetadata{}, bp.CollectMetadata("non-song"), "CollectMetadata did not return empty metadata for non-song")
	assert.Equal(t, SongMetadata{Title: "test-title", Artist: "test-artist", Album: "test-album"},
		bp.CollectMetadata("test-song.mp3"), "CollectMetadata did not return the correct metadata")
}
