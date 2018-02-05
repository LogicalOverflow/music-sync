package main

import (
	"github.com/LogicalOverflow/music-sync/schedule"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testSampleRate = 1

var removeOldChunksCases = []struct {
	now    int64
	chunks []upcomingChunk
	result []upcomingChunk
}{
	{
		now:    0e9,
		chunks: []upcomingChunk{},
		result: []upcomingChunk{},
	},
	{
		now:    0e9,
		chunks: []upcomingChunk{{startTime: 0e9, startIndex: 0, size: 16}, {startTime: 16e9, startIndex: 16, size: 16}},
		result: []upcomingChunk{{startTime: 0e9, startIndex: 0, size: 16}, {startTime: 16e9, startIndex: 16, size: 16}},
	},
	{
		now:    8e9,
		chunks: []upcomingChunk{{startTime: 0e9, startIndex: 0, size: 16}, {startTime: 16e9, startIndex: 16, size: 16}},
		result: []upcomingChunk{{startTime: 0e9, startIndex: 0, size: 16}, {startTime: 16e9, startIndex: 16, size: 16}},
	},
	{
		now:    24e9,
		chunks: []upcomingChunk{{startTime: 0e9, startIndex: 0, size: 16}, {startTime: 16e9, startIndex: 16, size: 16}},
		result: []upcomingChunk{{startTime: 16e9, startIndex: 16, size: 16}},
	},
	{
		now:    24e9,
		chunks: []upcomingChunk{{startTime: 0e9, startIndex: 0, size: 16}, {startTime: 32e9, startIndex: 32, size: 16}},
		result: []upcomingChunk{{startTime: 32e9, startIndex: 32, size: 16}},
	},
}

var currentSampleCases = []struct {
	now    int64
	chunks []upcomingChunk
	result int64
}{
	{
		now:    0,
		chunks: []upcomingChunk{},
		result: 0,
	},
	{
		now:    16e9,
		chunks: []upcomingChunk{},
		result: 0,
	},
	{
		now:    1e9,
		chunks: []upcomingChunk{{startTime: 16e9, startIndex: 16, size: 16}, {startTime: 32e9, startIndex: 32, size: 16}},
		result: 1,
	},
	{
		now:    8e9,
		chunks: []upcomingChunk{{startTime: 0e9, startIndex: 16, size: 16}, {startTime: 16, startIndex: 32, size: 16}},
		result: 24,
	},
	{
		now:    40e9,
		chunks: []upcomingChunk{{startTime: 16e9, startIndex: 32, size: 16}, {startTime: 32e9, startIndex: 48, size: 16}},
		result: 56,
	},
	{
		now:    40e9,
		chunks: []upcomingChunk{{startTime: 16e9, startIndex: 32, size: 16}, {startTime: 48e9, startIndex: 64, size: 16}},
		result: 56,
	},
}

var currentSongCases = []struct {
	sample int64
	songs  []upcomingSong
	result upcomingSong
}{
	{
		sample: 0,
		songs:  []upcomingSong{},
		result: upcomingSong{filename: "None", startIndex: 0, length: 0},
	},
	{
		sample: 0,
		songs:  []upcomingSong{{filename: "too-late", startIndex: 16, length: 16}},
		result: upcomingSong{filename: "None", startIndex: 0, length: 0},
	},
	{
		sample: 32,
		songs:  []upcomingSong{{filename: "passed", startIndex: 0, length: 16}, {filename: "current", startIndex: 24, length: 16}},
		result: upcomingSong{filename: "current", startIndex: 24, length: 16},
	},
	{
		sample: 8,
		songs:  []upcomingSong{{filename: "current", startIndex: 0, length: 16}, {filename: "next", startIndex: 24, length: 16}},
		result: upcomingSong{filename: "current", startIndex: 0, length: 16},
	},
}

var removeOldPausesCases = []struct {
	song   upcomingSong
	pauses []pauseToggle
	result []pauseToggle
}{
	{
		song:   upcomingSong{filename: "None", startIndex: 0, length: 0},
		pauses: []pauseToggle{},
		result: []pauseToggle{},
	},
	{
		song:   upcomingSong{filename: "None", startIndex: 0, length: 0},
		pauses: []pauseToggle{{playing: true, toggleIndex: 16}, {playing: false, toggleIndex: 32}},
		result: []pauseToggle{{playing: true, toggleIndex: 16}, {playing: false, toggleIndex: 32}},
	},
	{
		song:   upcomingSong{filename: "actual-song", startIndex: 64, length: 0},
		pauses: []pauseToggle{{playing: false, toggleIndex: 16}, {playing: true, toggleIndex: 32}, {playing: false, toggleIndex: 48}},
		result: []pauseToggle{{playing: true, toggleIndex: 32}, {playing: false, toggleIndex: 48}},
	},
	{
		song:   upcomingSong{filename: "actual-song", startIndex: 64, length: 0},
		pauses: []pauseToggle{{playing: false, toggleIndex: 32}, {playing: true, toggleIndex: 48}},
		result: []pauseToggle{{playing: true, toggleIndex: 48}},
	},
}

var pausesInCurrentSonCases = []struct {
	sample        int64
	song          upcomingSong
	pauses        []pauseToggle
	resultPauses  int64
	resultPlaying bool
}{
	{
		sample:        0,
		song:          upcomingSong{filename: "None", startIndex: 0, length: 0},
		pauses:        []pauseToggle{},
		resultPauses:  0,
		resultPlaying: true,
	},
	{
		sample:        128,
		song:          upcomingSong{filename: "actual-song", startIndex: 32, length: 256},
		pauses:        []pauseToggle{{playing: true, toggleIndex: 0}, {playing: false, toggleIndex: 16}, {playing: true, toggleIndex: 48}},
		resultPauses:  16,
		resultPlaying: true,
	},
	{
		sample:        128,
		song:          upcomingSong{filename: "actual-song", startIndex: 32, length: 256},
		pauses:        []pauseToggle{{playing: true, toggleIndex: 0}, {playing: false, toggleIndex: 64}},
		resultPauses:  64,
		resultPlaying: false,
	},
	{
		sample:        128,
		song:          upcomingSong{filename: "actual-song", startIndex: 16, length: 256},
		pauses:        []pauseToggle{{playing: true, toggleIndex: 0}, {playing: false, toggleIndex: 32}, {playing: true, toggleIndex: 64}},
		resultPauses:  32,
		resultPlaying: true,
	},
}

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
	for _, c := range pausesInCurrentSonCases {
		state := state{Pauses: c.pauses}
		actualPauses, actualPlaying := state.pausesInCurrentSong(c.sample, c.song)
		assert.Equal(t, c.resultPauses, actualPauses, "pausesInCurrentSong returned wrong pause for case %v", c)
		assert.Equal(t, c.resultPlaying, actualPlaying, "pausesInCurrentSong returned wrong playing for case %v", c)
	}
}
