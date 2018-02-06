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

func createSampleSlice(start, count int) [][2]float64 {
	s := make([][2]float64, count)
	for i := range s {
		s[i] = [2]float64{-float64(start + i), float64(start + i)}
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

func TestQueueChunk(t *testing.T) {
	var oldStreamer *timedMultiStreamer
	if streamer != nil {
		*oldStreamer = *streamer
	}
	streamer = &timedMultiStreamer{chunks: make([]*queuedChunk, 0)}

	startTime := make([]int64, 16)
	samples := make([][][2]float64, 16)
	for i := 0; i < 16; i++ {
		startTime[i] = 512 * int64(i)
		samples[i] = createSampleSlice(i*512, 512)
	}

	for i := 0; i < 16; i++ {
		QueueChunk(startTime[i], int64(i), samples[i])
		if !assert.Equal(t, i+1, len(streamer.chunks), "streamer.chunks has wrong length after queueing %d chunks", i+1) {
			continue
		}
		for j, qc := range streamer.chunks {
			assert.Equal(t, startTime[j], qc.startTime, "%d-th streamer.chunks has the wrong start time", j)
			assert.Equal(t, samples[j], qc.samples, "%d-th streamer.chunks has the wrong samples", j)
			assert.Equal(t, 0, qc.pos, "%d-th streamer.chunks has the wrong pos", j)
			assert.Equal(t, len(samples[j]), qc.sampleN, "%d-th streamer.chunks has the wrong sampleN", j)
		}
	}

	if oldStreamer != nil {
		*streamer = *oldStreamer
	}
}

func TestSamplesToAudioBuf(t *testing.T) {
	oldVolume := volume

	samples := createSampleSlice(0, 1024)
	buf := make([]byte, 4*len(samples))
	for volume = 0; volume <= 1; volume += 0.125 {
		samplesToAudioBuf(samples, buf)
		for i := 0; i < len(samples); i++ {
			ell, elh := convertSampleToBytes(samples[i][0] * volume)
			ehl, ehh := convertSampleToBytes(samples[i][1] * volume)
			assert.Equal(t, ell, buf[4*i+0], "samplesToAudioBuf has the wrong lower byte for the lower sample at index %d with volume %f", i, volume)
			assert.Equal(t, elh, buf[4*i+1], "samplesToAudioBuf has the wrong higher byte for the lower sample at index %d with volume %f", i, volume)
			assert.Equal(t, ehl, buf[4*i+2], "samplesToAudioBuf has the wrong lower byte for the higher sample at index %d with volume %f", i, volume)
			assert.Equal(t, ehh, buf[4*i+3], "samplesToAudioBuf has the wrong higher byte for the higher sample at index %d with volume %f", i, volume)
		}
	}

	volume = oldVolume
}

func TestConvertSampleToBytes(t *testing.T) {
	for i := -32768; i < 32767; i++ {
		el, eh := byte(i), byte(i>>8)
		f := float64(i) / float64(32768)
		l, h := convertSampleToBytes(f)
		assert.Equal(t, el, l, "convertSampleToBytes lower byte is wrong for i=%d", i)
		assert.Equal(t, eh, h, "convertSampleToBytes higher byte is wrong for i=%d", i)
	}
}
