package playback

import (
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/faiface/beep"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newSongsList(count int) []string {
	s := make([]string, count)
	for i := 0; i < count; i++ {
		s[i] = songName(i)
	}
	return s
}

func TestCombineSamples(t *testing.T) {
	testSamplesCount := 64
	low := make([]float64, testSamplesCount)
	high := make([]float64, testSamplesCount)

	for i := 0; i < testSamplesCount; i++ {
		low[i] = float64(i)
		high[i] = float64(testSamplesCount + i)
	}

	combined := CombineSamples(low, high)

	for i := 0; i < testSamplesCount; i++ {
		assert.Equal(t, float64(i), low[i], "CombineSamples modified low samples at index %d", i)
		assert.Equal(t, float64(testSamplesCount+i), high[i], "CombineSamples modified high samples at index %d", i)

		assert.Equal(t, low[i], combined[i][0], "low combined sample at index %d has the wrong value", i)
		assert.Equal(t, high[i], combined[i][1], "high combined sample at index %d has the wrong value", i)
	}
}

func TestSetVolume(t *testing.T) {
	log.DefaultCutoffLevel = log.LevelOff
	for f := float64(0); f <= 1; f += .1 {
		SetVolume(f)
		assert.Equal(t, f, volume, "SetVolume did not set volume correctly")
	}
}

func TestGetStreamer(t *testing.T) {
	log.DefaultCutoffLevel = log.LevelOff
	ad := AudioDir
	AudioDir = "_playback_test_files"

	var err error
	var s beep.StreamSeekCloser

	s, err = getStreamer("non-existent")
	assert.NotNil(t, err, "getStreamer did not return an error for a non-existent file")
	s, err = getStreamer("bad-format.mp3")
	assert.NotNil(t, err, "getStreamer did not return an error for a file with a bad format")
	s, err = getStreamer("okay.mp3")
	assert.Nil(t, err, "getStreamer did return an error for an okay file")
	assert.Equal(t, 443520, s.Len(), "getStreamer's stream returned an incorrect length")

	AudioDir = ad
}
