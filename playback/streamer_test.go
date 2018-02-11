package playback

import (
	"context"
	"github.com/faiface/beep"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueuedChunk(t *testing.T) {
	for i := 0; i < 16; i++ {
		samples := createSampleSlice(1024*i, 1024)

		qs := newQueuedStream(int64(i), samples)
		assert.Equal(t, 0, qs.pos, "the %d-th chunk is not initialized with position 0", i+1)
		assert.Equal(t, 1024, qs.sampleN, "the %d-th chunk is not initialized sampleN 1024", i+1)
		for j := 0; j < 7; j++ {
			result := make([][2]float64, 128)
			n := qs.copySamples(result)
			assert.Equal(t, 128, n, "%d-th copySample returned the wrong number of copied samples for chunk %d", j+1, i+1)
			assert.Equal(t, 128*(j+1), qs.pos, "%d-th copySample set the position wrong for the %d-th chunk", j+1, i+1)
			assert.Equal(t, samples[128*j:128*(j+1)], result, "%d-th copySample returned the wrong samples for the %d-th chunk", j+1, i+1)
			assert.False(t, qs.drained(), "after %d-th copySample, drained() returned true for the %d-th chunk", j+1, i+1)
		}

		result := make([][2]float64, 256)
		n := qs.copySamples(result)
		assert.Equal(t, 128, n, "copySample (under-filling result) returned the wrong number of copied samples for chunk %d", i+i)
		assert.Equal(t, 1024, qs.pos, "copySample (under-filling result) set the position wrong for the %d-th chunk", i+1)
		assert.Equal(t, samples[896:], result[:128], "copySample (under-filling result) returned the wrong samples for the %d-th chunk", i+1)
		assert.True(t, qs.drained(), "after copySample (under-filling result), drained() returned false for the %d-th chunk", i+1)

		n = qs.copySamples(result)
		assert.Equal(t, 0, n, "last copySample (after drained) returned the wrong number of copied samples for chunk %d", i+i)
		assert.Equal(t, 1024, qs.pos, "last copySample (after drained) set the position wrong for the %d-th chunk", i+1)
		assert.True(t, qs.drained(), "after last copySample (after drained), drained() returned false for the %d-th chunk", i+1)
	}
}

func TestTimedMultiStream_Err(t *testing.T) {
	tms := &timedMultiStreamer{}
	assert.Nil(t, tms.Err(), "timedMultiStreamer Err returned not nil")
}

func TestTimedMultiStream_sampleDurationAndSampleCount(t *testing.T) {
	tms := &timedMultiStreamer{format: beep.Format{SampleRate: 1000}}
	for i := 1; i < 1e6; i *= 10 {
		assert.Equal(t, int64(i*1e6), tms.samplesDuration(i), "timeDuration is wrong for 1 sample")
		assert.Equal(t, i, tms.samplesCount(int64(i*1e6)), "timeDuration is wrong for 1 sample")
	}
}

func TestTimedMultiStreamer_ReadChunks(t *testing.T) {
	tms := &timedMultiStreamer{
		format:  beep.Format{SampleRate: 1},
		chunks:  make([]*queuedChunk, 0),
		samples: newTimedSampleQueue(64),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go tms.ReadChunks(ctx)

	tms.chunksMutex.Lock()
	for i := 0; i < 16; i++ {
		tms.chunks = append(tms.chunks, newTestChunk(128, i))
	}
	tms.chunksMutex.Unlock()

	for i := 0; i < 2048; i++ {
		sample, time := tms.samples.Remove()
		assert.Equal(t, [2]float64{-float64(i), float64(i)}, sample, "timeDuration ReadChunks pushed wrong sample at index %d", i)
		assert.Equal(t, int64(i*1e9), time, "timeDuration ReadChunks pushed wrong times at index %d", i)
	}

	cancel()
}

func newTestChunk(chunkSize, chunkNum int) *queuedChunk {
	qc := &queuedChunk{
		startTime: int64(chunkSize * chunkNum * 1e9),
		samples:   createSampleSlice(chunkSize*chunkNum, chunkSize),
		sampleN:   chunkSize,
	}
	return qc
}
