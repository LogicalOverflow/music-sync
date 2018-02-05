package main

import (
	"github.com/LogicalOverflow/music-sync/schedule"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestState_removeOldChunks(t *testing.T) {
	schedule.SampleRate = testSampleRate
	for _, c := range removeOldChunksCases {
		state := state{Chunks: c.chunks}
		state.removeOldChunks(c.now)
		if !assert.Equal(t, len(c.result), len(state.Chunks), "removeOldChunks removed the wrong number of chunks for case %v: %v", c, state.Chunks) {
			continue
		}
		for i := range c.result {
			assert.Equal(t, c.result[i], state.Chunks[i], "removeOldChunks resulted with the wrong chunk at index %d for case %v", i, c)
		}
	}

}

func TestState_currentSample(t *testing.T) {
	schedule.SampleRate = testSampleRate
	for _, c := range currentSampleCases {
		state := state{Chunks: c.chunks}
		actual := state.currentSample(c.now)
		assert.Equal(t, c.result, actual, "currentSample is incorrect at %d with chunks %v", c.now, c.chunks)
	}
}

func TestState_currentSong(t *testing.T) {
	schedule.SampleRate = testSampleRate
	for _, c := range currentSongCases {
		state := state{Songs: c.songs}
		actual := state.currentSong(c.sample)
		assert.Equal(t, c.result, actual, "currentSong is incorrect at sample %d with songs %v", c.sample, c.songs)
	}
}

func TestState_removeOldPauses(t *testing.T) {
	schedule.SampleRate = testSampleRate
	for _, c := range removeOldPausesCases {
		state := state{Pauses: c.pauses}
		state.removeOldPauses(c.song)
		if !assert.Equal(t, len(c.result), len(state.Pauses), "removeOldPauses removed the wrong number of pauses for case %v: %v", c, state.Pauses) {
			continue
		}
		for i := range c.result {
			assert.Equal(t, c.result[i], state.Pauses[i], "removeOldPauses resulted with the wrong pause at index %d for case %v", i, c)
		}
	}
}

func TestState_pausesInCurrentSong(t *testing.T) {
	schedule.SampleRate = testSampleRate
	for _, c := range pausesInCurrentSongCases {
		state := state{Pauses: c.pauses}
		actualPauses, actualPlaying := state.pausesInCurrentSong(c.sample, c.song)
		assert.Equal(t, c.resultPauses, actualPauses, "pausesInCurrentSong returned wrong pause for case %v", c)
		assert.Equal(t, c.resultPlaying, actualPlaying, "pausesInCurrentSong returned wrong playing for case %v", c)
	}
}

func TestState_Info(t *testing.T) {
	for _, c := range infoTestCases {
		actual := c.state.Info(c.now)
		assert.Equal(t, c.info, *actual, "state Info returned wrong info for case %v", c)
	}
}

func TestPlaybackInformation_playingString(t *testing.T) {
	assert.Equal(t, "Playing", playbackInformation{Playing: true}.playingString(), "playbackString is wrong for Playing: true")
	assert.Equal(t, "Paused", playbackInformation{Playing: false}.playingString(), "playbackString is wrong for Playing: false")
}

func TestUpcomingChunk_lengthAndEndTime(t *testing.T) {
	schedule.SampleRate = testSampleRate
	for _, c := range upcomingChunkCases {
		assert.Equal(t, c.length, c.chunk.length(), "upcomingChunk %v has wrong length", c.chunk)
		assert.Equal(t, c.endTime, c.chunk.endTime(), "upcomingChunk %v has wrong end time", c.chunk)
	}
}
