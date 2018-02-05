package main

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

var pausesInCurrentSongCases = []struct {
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

var infoTestCases = []struct {
	now   int64
	state state
	info  playbackInformation
}{
	{
		now: 0e9,
		state: state{
			Songs:  []upcomingSong{},
			Chunks: []upcomingChunk{},
			Pauses: []pauseToggle{},
			Volume: 0.1,
		},
		info: playbackInformation{
			CurrentSong:         upcomingSong{filename: "None", startIndex: 0, length: 0},
			CurrentSample:       0,
			PausesInCurrentSong: 0,
			Now:                 0e9,
			Playing:             true,
			Volume:              0.1,
			SongLength:          0,
			TimeInSong:          0,
			ProgressInSong:      0,
		},
	},
	{
		now: 45e9,
		state: state{
			Songs:  []upcomingSong{{filename: "the-song", startIndex: 15, length: 60}},
			Chunks: []upcomingChunk{{startTime: 0e9, startIndex: 0, size: 256}},
			Pauses: []pauseToggle{},
			Volume: 0.1,
		},
		info: playbackInformation{
			CurrentSong:         upcomingSong{filename: "the-song", startIndex: 15, length: 60},
			CurrentSample:       45,
			PausesInCurrentSong: 0,
			Now:                 45e9,
			Playing:             true,
			Volume:              0.1,
			SongLength:          60e9,
			TimeInSong:          30e9,
			ProgressInSong:      0.5,
		},
	},
	{
		now: 60e9,
		state: state{
			Songs:  []upcomingSong{{filename: "the-song", startIndex: 15, length: 60}},
			Chunks: []upcomingChunk{{startTime: 0e9, startIndex: 0, size: 256}},
			Pauses: []pauseToggle{{playing: true, toggleIndex: 0}, {playing: false, toggleIndex: 30}, {playing: true, toggleIndex: 45}},
			Volume: 0.1,
		},
		info: playbackInformation{
			CurrentSong:         upcomingSong{filename: "the-song", startIndex: 15, length: 60},
			CurrentSample:       60,
			PausesInCurrentSong: 15,
			Now:                 60e9,
			Playing:             true,
			Volume:              0.1,
			SongLength:          60e9,
			TimeInSong:          30e9,
			ProgressInSong:      0.5,
		},
	},
}

var upcomingChunkCases = []struct {
	chunk   upcomingChunk
	length  int64
	endTime int64
}{
	{
		chunk:   upcomingChunk{startTime: 0e9, startIndex: 0, size: 0},
		length:  0,
		endTime: 0,
	},
	{
		chunk:   upcomingChunk{startTime: 0e9, startIndex: 16, size: 16},
		length:  16e9,
		endTime: 16e9,
	},
	{
		chunk:   upcomingChunk{startTime: 16e9, startIndex: 16, size: 16},
		length:  16e9,
		endTime: 32e9,
	},
}
