package playback

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	if assert.Equal(t, 16-len(removed), len(pl.songs), "after removing %v, playlist songs has the wrong length", removed) {
		skipped := 0
		for i := 0; i < 16; i++ {
			skip := false
			for _, r := range removed {
				if r == i {
					skip = true
				}
			}

			if skip {
				skipped++
			} else {
				assert.Equal(t, songName(i), pl.songs[i-skipped], "after removing %v, playlist song at index %d is incorrect", i-skipped)
			}
		}
	}
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
