package playback

import (
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
