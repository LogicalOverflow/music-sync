package playback

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPlaylist_pushStreamer(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	pl.currentSong = "the-song"
	pl.sampleIndexWrite = 16

	msg := "%d-th %s sample is incorrect when using pushStreamer"
	msgNan := "%d-th %s pause sample is not none when using pushStreamer"

	chunkCount := 2

	var startSampleIndex uint64
	var filename string
	var songLength int64

	pl.SetNewSongHandler(storingNewSongHandler(&startSampleIndex, &filename, &songLength))

	s := &testStreamer{samples: make(chan [2]float64, 512), position: 0, length: 1024}

	pl.SetPlaying(true)
	go pl.pushStreamer(s)

	comm := make(chan bool)

	go s.pushSamplesInChunksWithPausesAndClose(pl, streamerBufferSize, chunkCount, comm)

	for c := 0; c < chunkCount; c++ {
		assertPlaylistSamplesInChan(t, pl, c*streamerBufferSize, 1, msg)
		comm <- true
		assertPlaylistSamplesInChan(t, pl, c*streamerBufferSize+1, streamerBufferSize-1, msg)

		assertPlaylistNanSamplesInChan(t, pl, c*streamerBufferSize, 1, msgNan)
		comm <- true
		assertPlaylistNanSamplesInChan(t, pl, c*streamerBufferSize+1, streamerBufferSize-1, msgNan)
	}

	assertPlaylistSamplesInChan(t, pl, chunkCount*streamerBufferSize, streamerBufferSize, msg)

	assert.Equal(t, uint64(16), startSampleIndex, "NewSongHandler called with wrong startSampleIndex")
	assert.Equal(t, "the-song", filename, "NewSongHandler called with wrong filename")
	assert.Equal(t, int64(1024), songLength, "NewSongHandler called with wrong songLength")
}

func TestPlaylist_pushSample(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	go func() {
		for i := 0; i < 32; i++ {
			pl.pushSample(-float64(i), float64(i))
		}
	}()
	for i := 0; i < 32; i++ {
		assert.Equal(t, -float64(i), <-pl.low, "%d-th low sample is wrong when pushing with pushSample", i)
		assert.Equal(t, float64(i), <-pl.high, "%d-th high sample is wrong when pushing with pushSample", i)
	}
}

func TestPlaylist_pushNanSamples(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	go pl.pushNanSamples(32)
	for i := 0; i < 32; i++ {
		assert.True(t, math.IsNaN(<-pl.low), "%d-th low sample is not nan when pushNanSamples", i)
		assert.True(t, math.IsNaN(<-pl.high), "%d-th high sample is not nan when pushNanSamples", i)
	}
}

func TestPlaylist_pushBuffer(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	buffer := make([][2]float64, 32)
	for i := range buffer {
		buffer[i] = [2]float64{-float64(i), float64(i)}
	}
	go pl.pushBuffer(buffer)

	for i := range buffer {
		assert.Equal(t, buffer[i][0], <-pl.low, "%d-th low sample is wrong when pushing with pushBuffer", i)
		assert.Equal(t, buffer[i][1], <-pl.high, "%d-th high sample is wrong when pushing with pushBuffer", i)
	}
}

func TestPlaylist_pushNanBreak(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 32)
	go pl.pushNanBreak()
	for i := 0; i < 32; i++ {
		assert.True(t, math.IsNaN(<-pl.low), "%d-th low sample is not nan when pushNanBreak", i)
		assert.True(t, math.IsNaN(<-pl.high), "%d-th high sample is not nan when pushNanBreak", i)
	}
}
