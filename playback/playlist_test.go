package playback

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"sync"
	"testing"
)

func songName(index int) string {
	return fmt.Sprintf("song-%02d", index)
}

func TestPlaylist_SetPos(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	for i := 0; i < 16; i++ {
		go func() { <-pl.forceNext }()
		pl.SetPos(i)
		assert.Equal(t, i, pl.position, "playlist SetPos did not set position correctly")
	}
}

func TestPlaylist_Songs(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	for i := 0; i < 16; i++ {
		pl.songs = make([]string, i+1)
		for j := range pl.songs {
			pl.songs[j] = songName(i)
		}
		songs := pl.Songs()
		if assert.Equal(t, i+1, len(songs), "playlist Songs returned a slice of incorrect length when holding %d songs", i+1) {
			for j := range songs {
				assert.Equal(t, songName(i), songs[j], "playlist Songs returned a slice with the wrong song at index %d when holding %d songs", j, i+1)
			}
		}
	}
}

func TestPlaylist_AddSong(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	require.Zero(t, len(pl.songs), "newly created playlist contains songs")
	for i := 0; i < 16; i++ {
		pl.AddSong(songName(i))
		if assert.Equal(t, i+1, len(pl.songs), "after adding %d songs, playlist does not hold the right amount of songs", i+1) {
			for j := 0; j <= i; j++ {
				assert.Equal(t, songName(i), pl.songs[i], "after adding %d songs, the song with index %d is incorrect", i+1, j)
			}
		}
	}
}

func TestPlaylist_InsertSong(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	pl.songs = make([]string, 8)
	for i := range pl.songs {
		pl.songs[i] = songName(2 * i)
	}

	for i := 1; i < 16; i += 2 {
		pl.InsertSong(songName(i), i)
		if assert.Equal(t, i/2+9, len(pl.songs), "after inserting %d test songs, songs length is incorrect", i/2+1) {
			for j := 0; j <= i; j++ {
				assert.Equal(t, songName(j), pl.songs[j], "after inserting %d test songs, song at index %d is incorrect", i/2+1, j)
			}
		}
	}

	expectedSongs := make([]string, 16)
	for i := range expectedSongs {
		expectedSongs[i] = songName(i)
	}

	assert.Equal(t, expectedSongs, pl.songs, "after inserting 8 songs, songs is incorrect")

	pl.InsertSong("song-low", -1)
	expectedSongs = append([]string{"song-low"}, expectedSongs...)
	assert.Equal(t, expectedSongs, pl.songs, "after inserting 9 songs, songs is incorrect")

	pl.InsertSong("song-high", 32)
	expectedSongs = append(expectedSongs, "song-high")
	assert.Equal(t, expectedSongs, pl.songs, "after inserting 10 songs, songs is incorrect")
}

func TestPlaylist_RemoveSong(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	pl.songs = make([]string, 16)
	for i := range pl.songs {
		pl.songs[i] = songName(i)
	}

	pl.RemoveSong(8)
	assertRemoved(t, []int{8}, pl)

	pl.RemoveSong(10)
	assertRemoved(t, []int{8, 11}, pl)

	pl.RemoveSong(1)
	assertRemoved(t, []int{1, 8, 11}, pl)

	pl.RemoveSong(-2)
	assertRemoved(t, []int{0, 1, 8, 11}, pl)

	pl.RemoveSong(22)
	assertRemoved(t, []int{0, 1, 8, 11, 15}, pl)
}

func assertRemoved(t *testing.T, removed []int, pl *Playlist) {
	expected := make([]string, 16-len(removed))
	skipped := 0
	for i := 0; i < 16; i++ {
		if intSliceContains(removed, i) {
			skipped++
		} else {
			expected[i-skipped] = songName(i)
		}
	}

	assert.Equal(t, expected, pl.songs, "after removing %v, playlist songs are incorrect", removed)
}

func intSliceContains(ints []int, t int) bool {
	for _, i := range ints {
		if t == i {
			return true
		}
	}
	return false
}

func TestPlaylist_Fill(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)

	in := make([][2]float64, 64)
	for i := range in {
		in[i][0] = float64(i)
		in[i][1] = -float64(i)
	}
	go func() {
		for _, s := range in {
			pl.low <- s[0]
			pl.high <- s[1]
		}
	}()

	low := make([]float64, 64)
	high := make([]float64, 64)

	pl.Fill(low, high)

	for i := range in {
		assert.Equal(t, in[i][0], low[i], "playlist fill inserted the wrong low sample at index %d", i)
		assert.Equal(t, in[i][1], high[i], "playlist fill inserted the wrong high sample at index %d", i)
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

func TestPlaylist_pushStreamer(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	pl.currentSong = "the-song"
	pl.sampleIndexWrite = 16

	var startSampleIndex uint64
	var filename string
	var songLength int64

	pl.SetNewSongHandler(func(ssi uint64, fn string, sl int64) {
		startSampleIndex = ssi
		filename = fn
		songLength = sl
	})

	s := &testStreamer{samples: make(chan [2]float64, 512), position: 0, length: 1024}

	pl.SetPlaying(true)
	go pl.pushStreamer(s)

	comm := make(chan bool)

	go func() {
		for i := 0; i < 512; i++ {
			s.samples <- [2]float64{-float64(i), float64(i)}
		}
		<-comm
		pl.SetPlaying(false)
		<-comm
		pl.SetPlaying(true)
		for i := 512; i < 1024; i++ {
			s.samples <- [2]float64{-float64(i), float64(i)}
		}
		s.Close()
	}()

	assert.Equal(t, -float64(0), <-pl.low, "0-th low sample is incorrect when using pushStreamer")
	assert.Equal(t, float64(0), <-pl.high, "0-th high sample is incorrect when using pushStreamer")
	comm <- true
	for i := 1; i < 512; i++ {
		assert.Equal(t, -float64(i), <-pl.low, "%d-th low sample is incorrect when using pushStreamer", i)
		assert.Equal(t, float64(i), <-pl.high, "%d-th high sample is incorrect when using pushStreamer", i)
	}

	assert.True(t, math.IsNaN(<-pl.low), "0-th low pause sample is not none when using pushStreamer")
	assert.True(t, math.IsNaN(<-pl.high), "0-th high pause sample is not nan when using pushStreamer")
	comm <- true
	for i := 1; i < 512; i++ {
		assert.True(t, math.IsNaN(<-pl.low), "%d-th low pause sample is not none when using pushStreamer", i)
		assert.True(t, math.IsNaN(<-pl.high), "%d-th high pause sample is not nan when using pushStreamer", i)
	}
	for i := 512; i < 1024; i++ {
		assert.Equal(t, -float64(i), <-pl.low, "%d-th low sample is incorrect when using pushStreamer", i)
		assert.Equal(t, float64(i), <-pl.high, "%d-th high sample is incorrect when using pushStreamer", i)
	}

	assert.Equal(t, uint64(16), startSampleIndex, "NewSongHandler called with wrong startSampleIndex")
	assert.Equal(t, "the-song", filename, "NewSongHandler called with wrong filename")
	assert.Equal(t, int64(1024), songLength, "NewSongHandler called with wrong songLength")
}

func TestPlaylist_shouldBreakStreamerPushLoop(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	assert.True(t, pl.shouldBreakStreamerPushLoop(16, false, 16), "playlist shouldBreakStreamerPushLoop returned false when ok is false")
	assert.True(t, pl.shouldBreakStreamerPushLoop(15, true, 16), "playlist shouldBreakStreamerPushLoop returned false when n < bufSize")
	assert.False(t, pl.shouldBreakStreamerPushLoop(16, true, 16), "playlist shouldBreakStreamerPushLoop returned true when n = bufSize and ok true")
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

func TestPlaylist_callPauseToggleHandler(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	var playing bool
	var sample uint64
	var wg sync.WaitGroup
	pl.SetPauseToggleHandler(func(p bool, s uint64) {
		wg.Done()
		playing = p
		sample = s
	})

	for i := 0; i < 4; i++ {
		sample = uint64(0xffffffffffffffff)
		pl.sampleIndexWrite = uint64(i)
		plPlayingLast := i&1 == 1
		plPlaying := i&2 == 2

		pl.playingLast = plPlayingLast
		pl.playing = plPlaying
		if plPlaying != plPlayingLast {
			wg.Add(1)
		}

		pl.callPauseToggleHandler()
		if plPlaying == plPlayingLast {
			assert.Equal(t, uint64(0xffffffffffffffff), sample, "PauseToggleHandler called with playing %v and playingLast %v", plPlaying, plPlayingLast)
		} else {
			wg.Wait()
			assert.Equal(t, uint64(i), sample, "PauseToggleHandler not called with correct sample for playing %v and playingLast %v", plPlaying, plPlayingLast)
			assert.Equal(t, plPlaying, playing, "PauseToggleHandler not called with correct playing for playing %v and playingLast %v", plPlaying, plPlayingLast)
		}
	}
	pl.playingLast = true
	pl.playing = true
}
