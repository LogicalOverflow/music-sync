package main

import (
	"github.com/LogicalOverflow/music-sync/metadata"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testLyrics = []metadata.LyricsLine{
	{{Timestamp: 10000, Caption: "caption-line-1-caption-0 "}, {Timestamp: 11000, Caption: "caption-line-1-caption-1 "}, {Timestamp: 12000, Caption: "caption-line-1-caption-2"}},
	{{Timestamp: 20000, Caption: "caption-line-2-caption-0 "}, {Timestamp: 21000, Caption: "caption-line-2-caption-1 "}, {Timestamp: 22000, Caption: "caption-line-2-caption-2"}},
	{{Timestamp: 30000, Caption: "caption-line-3-caption-0 "}, {Timestamp: 31000, Caption: "caption-line-3-caption-1 "}, {Timestamp: 32000, Caption: "caption-line-3-caption-2"}},
}

var lyricsBuildLineCases = []struct {
	index      int
	lyrics     []metadata.LyricsLine
	timeInSong time.Duration
	result     string
}{
	{
		index:      -1,
		lyrics:     testLyrics,
		timeInSong: 100 * time.Second,
		result:     "",
	},
	{
		index:      99,
		lyrics:     testLyrics,
		timeInSong: 100 * time.Second,
		result:     "",
	},
	{
		index:      0,
		lyrics:     testLyrics,
		timeInSong: 11 * time.Second,
		result:     "caption-line-1-caption-0 caption-line-1-caption-1 ",
	},
	{
		index:      0,
		lyrics:     testLyrics,
		timeInSong: 21 * time.Second,
		result:     "caption-line-1-caption-0 caption-line-1-caption-1 caption-line-1-caption-2",
	},
	{
		index:      0,
		lyrics:     testLyrics,
		timeInSong: 0,
		result:     "",
	},
	{
		index:      1,
		lyrics:     testLyrics,
		timeInSong: 100 * time.Second,
		result:     "caption-line-2-caption-0 caption-line-2-caption-1 caption-line-2-caption-2",
	},
}

func TestLyricsBuildLine(t *testing.T) {
	for _, c := range lyricsBuildLineCases {
		s := upcomingSong{lyrics: c.lyrics}
		actual := lyricsBuildLine(c.index, s, c.timeInSong)
		assert.Equal(t, c.result, actual, "lyricsBuildLine returned an incorrect line for case %v", c)
	}
}

var lyricsNextLineCases = []struct {
	lyrics     []metadata.LyricsLine
	timeInSong time.Duration
	result     int
}{
	{
		lyrics:     []metadata.LyricsLine{},
		timeInSong: 100 * time.Second,
		result:     0,
	},
	{
		lyrics:     testLyrics,
		timeInSong: 0,
		result:     0,
	},
	{
		lyrics:     testLyrics,
		timeInSong: 5 * time.Second,
		result:     0,
	},
	{
		lyrics:     testLyrics,
		timeInSong: 10 * time.Second,
		result:     1,
	},
	{
		lyrics:     testLyrics,
		timeInSong: 20 * time.Second,
		result:     2,
	},
}

func TestLyricsNextLine(t *testing.T) {
	for _, c := range lyricsNextLineCases {
		s := upcomingSong{lyrics: c.lyrics}
		actual := lyricsNextLine(s, c.timeInSong)
		assert.Equal(t, c.result, actual, "lyricsNextLine result is wrong for case %v", c)
	}
}

var lyricsHistoryCases = []struct {
	lyrics     []metadata.LyricsLine
	timeInSong time.Duration
	histSize   int
	result     []string
}{
	{
		lyrics:     testLyrics,
		timeInSong: 30 * time.Second,
		histSize:   2,
		result:     []string{"caption-line-3-caption-0 ", "caption-line-2-caption-0 caption-line-2-caption-1 caption-line-2-caption-2"},
	},
	{
		lyrics:     testLyrics,
		timeInSong: 0,
		histSize:   4,
		result:     []string{"", "", "", ""},
	},
	{
		lyrics:     []metadata.LyricsLine{},
		timeInSong: 100 * time.Second,
		histSize:   4,
		result:     []string{"", "", "", ""},
	},
	{
		lyrics:     testLyrics,
		timeInSong: 100 * time.Second,
		histSize:   4,
		result:     []string{"caption-line-3-caption-0 caption-line-3-caption-1 caption-line-3-caption-2", "caption-line-2-caption-0 caption-line-2-caption-1 caption-line-2-caption-2", "caption-line-1-caption-0 caption-line-1-caption-1 caption-line-1-caption-2", ""},
	},
}

func TestLyricsHistory(t *testing.T) {
	for _, c := range lyricsHistoryCases {
		s := upcomingSong{lyrics: c.lyrics}
		actual := lyricsHistory(c.histSize, s, c.timeInSong)
		assert.Equal(t, c.result, actual, "lyricsHistory is wrong for case %v", c)
	}
}
