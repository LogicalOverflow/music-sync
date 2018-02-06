package playback

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func storingNewSongHandler(startSampleIndex *uint64, filename *string, songLength *int64) func(uint64, string, int64) {
	return func(ssi uint64, fn string, sl int64) {
		*startSampleIndex = ssi
		*filename = fn
		*songLength = sl
	}
}

var newPlaylistCases = []struct {
	bufferSize   int
	songs        []string
	nanBreakSize int
}{
	{
		bufferSize:   16,
		songs:        []string{},
		nanBreakSize: 0,
	},
	{
		bufferSize:   1,
		songs:        []string{"song-1", "song-2", "song-3"},
		nanBreakSize: 0,
	},
	{
		bufferSize:   1,
		songs:        []string{},
		nanBreakSize: 48,
	},
	{
		bufferSize:   16,
		songs:        []string{"song-1", "song-2", "song-3"},
		nanBreakSize: 48,
	},
}

func songName(index int) string {
	return fmt.Sprintf("song-%02d", index)
}

func assertPlaylistSamplesInChan(t *testing.T, pl *Playlist, start, count int, message string) {
	for i := 0; i < count; i++ {
		assert.Equal(t, -float64(start+i), <-pl.low, message, start+i, "low")
		assert.Equal(t, +float64(start+i), <-pl.high, message, start+i, "high")
	}
}

func assertPlaylistNanSamplesInChan(t *testing.T, pl *Playlist, start, count int, message string) {
	for i := 0; i < count; i++ {
		assert.True(t, math.IsNaN(<-pl.low), message, start+i, "low")
		assert.True(t, math.IsNaN(<-pl.high), message, start+i, "high")
	}
}

type testStreamer struct {
	samples  chan [2]float64
	position int
	length   int
}

func (ts *testStreamer) Err() error       { return nil }
func (ts *testStreamer) Len() int         { return ts.length }
func (ts *testStreamer) Position() int    { return ts.position }
func (ts *testStreamer) Seek(p int) error { ts.position = p; return nil }
func (ts *testStreamer) Close() error     { close(ts.samples); return nil }

func (ts *testStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n = 0
	ok = true
	for s := range ts.samples {
		samples[n] = s
		n++
		if len(samples) <= n {
			break
		}
	}
	ts.position += n
	return
}

func (ts *testStreamer) pushSamples(start, count int) {
	for i := 0; i < count; i++ {
		ts.samples <- [2]float64{-float64(start + i), float64(start + i)}
	}
}

func (ts *testStreamer) pushSamplesInChunksWithPausesAndClose(pl *Playlist, chunkSize, chunkCount int, comm chan bool) {
	for c := 0; c < chunkCount; c++ {
		ts.pushSamples(chunkSize*c, chunkSize)
		<-comm
		pl.SetPlaying(false)
		<-comm
		pl.SetPlaying(true)
	}
	ts.pushSamples(chunkSize*chunkCount, chunkSize)
	ts.Close()
}
