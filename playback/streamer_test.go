package playback

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueuedChunk_copySample(t *testing.T) {
	for i := 0; i < 16; i++ {
		samples := createSampleSlice(1024*i, 1024)

		qs := newQueuedStream(int64(i), samples)
		assert.Equal(t, 0, qs.pos, "the %d-th chunk is not initialized with position 0", i+1)
		assert.Equal(t, 1024, qs.sampleN, "the %d-th chunk is not initialized sampleN 1024", i+1)
		for j := 0; j < 7; j++ {
			result := make([][2]float64, 128)
			qs.copySamples(result)
			assert.Equal(t, 128*(j+1), qs.pos, "%d-th copySample set the position wrong for the %d-th chunk", i+1, j+1)
			assert.Equal(t, samples[128*j:128*(j+1)], result, "%d-th copySample returned the wrong samples for the %d-th chunk", i+1, j+1)
		}

		result := make([][2]float64, 256)
		qs.copySamples(result)
		assert.Equal(t, 1024, qs.pos, "last copySample (under-filling result) set the position wrong for the %d-th chunk", i+1)
		assert.Equal(t, samples[896:], result[:128], "last copySample (under-filling result) returned the wrong samples for the %d-th chunk", i+1)
	}
}
