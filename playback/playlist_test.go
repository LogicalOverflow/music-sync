package playback

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestPlaylist_SetPos(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	for i := 0; i < 16; i++ {
		pl.SetPos(i)
		assert.Equal(t, i, pl.position, "playlist SetPos did not set position correctly")
		assert.True(t, <-pl.forceNext, "playlist SetPos did not set forceNext correctly")
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

	expectedSongs := newSongsList(16)

	assert.Equal(t, expectedSongs, pl.songs, "after inserting 8 songs, songs is incorrect")

	pl.InsertSong("song-low", -1)
	expectedSongs = append([]string{"song-low"}, expectedSongs...)
	assert.Equal(t, expectedSongs, pl.songs, "after inserting 9 songs, songs is incorrect")

	pl.InsertSong("song-high", 32)
	expectedSongs = append(expectedSongs, "song-high")
	assert.Equal(t, expectedSongs, pl.songs, "after inserting 10 songs, songs is incorrect")
}

func TestPlaylist_RemoveSong(t *testing.T) {
	pl := NewPlaylist(16, newSongsList(16), 0)

	assert.Equal(t, songName(8), pl.RemoveSong(8), "remove returned the wrong song name")
	assertRemoved(t, []int{8}, pl)

	assert.Equal(t, songName(11), pl.RemoveSong(10), "remove returned the wrong song name")
	assertRemoved(t, []int{8, 11}, pl)

	assert.Equal(t, songName(1), pl.RemoveSong(1), "remove returned the wrong song name")
	assertRemoved(t, []int{1, 8, 11}, pl)

	assert.Equal(t, songName(0), pl.RemoveSong(-2), "remove returned the wrong song name")
	assertRemoved(t, []int{0, 1, 8, 11}, pl)

	assert.Equal(t, songName(15), pl.RemoveSong(22), "remove returned the wrong song name")
	assertRemoved(t, []int{0, 1, 8, 11, 15}, pl)

	pl.songs = []string{}
	assert.Equal(t, "", pl.RemoveSong(0), "remove returned the wrong song name for playlist without songs")
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

func TestPlaylist_shouldBreakStreamerPushLoop(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	assert.True(t, pl.shouldBreakStreamerPushLoop(16, false, 16), "playlist shouldBreakStreamerPushLoop returned false when ok is false")
	assert.True(t, pl.shouldBreakStreamerPushLoop(15, true, 16), "playlist shouldBreakStreamerPushLoop returned false when n < bufSize")
	assert.False(t, pl.shouldBreakStreamerPushLoop(16, true, 16), "playlist shouldBreakStreamerPushLoop returned true when n = bufSize and ok true")
	pl.forceNext <- true
	assert.True(t, pl.shouldBreakStreamerPushLoop(16, true, 16), "playlist shouldBreakStreamerPushLoop returned false when forceNext")
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

func TestNewPlaylist(t *testing.T) {
	for _, c := range newPlaylistCases {
		pl := NewPlaylist(c.bufferSize, c.songs, c.nanBreakSize)
		assert.Equal(t, c.bufferSize, cap(pl.low), "playlist low chan has wrong capacity for case %v", c)
		assert.Equal(t, c.bufferSize, cap(pl.high), "playlist high chan has wrong capacity for case %v", c)
		assert.Equal(t, c.songs, pl.songs, "playlist has wrong songs for case %v", c)
		assert.Equal(t, c.nanBreakSize, pl.nanBreakSize, "playlist has wrong nanBreakSize for case %v", c)
	}
}

func TestPlaylist_Pos(t *testing.T) {
	pl := NewPlaylist(16, []string{}, 0)
	for j := 0; j < 16; j++ {
		pl.position = j
		assert.Equal(t, 0, pl.Pos(), "playlist pos returned the wrong value with 0 songs at position %d", j)
	}
	for i := 1; i < 16; i++ {
		pl := NewPlaylist(16, newSongsList(i), 0)
		for j := 0; j < 16; j++ {
			pl.position = j
			assert.Equal(t, j%i, pl.Pos(), "playlist pos returned the wrong value with %d songs at position %d", i, j)
		}
	}
}

func TestPlaylist_nextSong(t *testing.T) {
	pl := NewPlaylist(16, newSongsList(16), 0)
	for i := 1; i <= 64; i++ {
		pl.position++
		assert.Equal(t, songName(i%16), pl.nextSong(), "nextSong returned the wrong song name after calling it %d times", i)
		assert.Equal(t, i%16, pl.position, "nextSong returned set the position incorrectly after calling it %d times", i)
	}

	pl.songs = []string{}
	assert.Equal(t, "", pl.nextSong(), "nextSong returned the wrong song name for playlist with no songs")
}
