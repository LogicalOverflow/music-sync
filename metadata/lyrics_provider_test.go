package metadata

import (
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLyricsProvider(t *testing.T) {
	assert.NotNil(t, GetLyricsProvider(), "GetLyricsProvider returned nil")
}

func TestBasicLyricsProvider_CollectLyrics(t *testing.T) {
	playback.AudioDir = "_test_files"
	lp := basicLyricsProvider{}

	l := lp.CollectLyrics("test-song.mp3")
	assert.Equal(t, []LyricsLine{
		{{Timestamp: 10000, Caption: "caption-0-0"}, {Timestamp: 11000, Caption: "caption-0-1"}, {Timestamp: 12000, Caption: "caption-0-2"}, {Timestamp: 13000, Caption: "caption-0-3"}, {Timestamp: 14000, Caption: "caption-0-4"}},
		{{Timestamp: 20000, Caption: "caption-1-0"}, {Timestamp: 21000, Caption: "caption-1-1"}, {Timestamp: 22000, Caption: "caption-1-2"}, {Timestamp: 23000, Caption: "caption-1-3"}, {Timestamp: 24000, Caption: "caption-1-4"}},
		{{Timestamp: 30000, Caption: "caption-2-0"}, {Timestamp: 31000, Caption: "caption-2-1"}, {Timestamp: 32000, Caption: "caption-2-2"}, {Timestamp: 33000, Caption: "caption-2-3"}, {Timestamp: 34000, Caption: "caption-2-4"}},
	}, l, "CollectLyrics did not return the correct lyrics")
	assert.Zero(t, len(lp.CollectLyrics("non-song")), "CollectLyrics did not return an empty line slice for a non-song")
}
