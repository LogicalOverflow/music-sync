package main

import (
	"github.com/LogicalOverflow/music-sync/metadata"
	"github.com/LogicalOverflow/music-sync/schedule"
	"sync"
	"time"
)

type state struct {
	Songs      []upcomingSong
	SongsMutex sync.RWMutex

	Chunks      []upcomingChunk
	ChunksMutex sync.RWMutex

	Pauses      []pauseToggle
	PausesMutex sync.RWMutex

	Volume float64
}

type pauseByToggleIndex []pauseToggle

func (p pauseByToggleIndex) Len() int           { return len(p) }
func (p pauseByToggleIndex) Less(i, j int) bool { return p[i].toggleIndex < p[j].toggleIndex }
func (p pauseByToggleIndex) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type chunksByStartIndex []upcomingChunk

func (c chunksByStartIndex) Len() int           { return len(c) }
func (c chunksByStartIndex) Less(i, j int) bool { return c[i].startIndex < c[j].startIndex }
func (c chunksByStartIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type songsByStartIndex []upcomingSong

func (s songsByStartIndex) Len() int           { return len(s) }
func (s songsByStartIndex) Less(i, j int) bool { return s[i].startIndex < s[j].startIndex }
func (s songsByStartIndex) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *state) Info(now int64) *playbackInformation {
	sample := s.currentSample(now)
	currentSong := s.currentSong(sample)
	s.removeOldPauses(currentSong)
	pausesInCurrentSong, playing := s.pausesInCurrentSong(sample, currentSong)

	var songLength, timeInSong time.Duration
	var progressInSong float64

	if currentSong.startIndex != 0 && int64(currentSong.startIndex) < sample {
		sampleInSong := sample - int64(currentSong.startIndex) - pausesInCurrentSong
		timeInSong = time.Duration(sampleInSong) * time.Second / time.Duration(schedule.SampleRate) / time.Nanosecond
		if 0 < currentSong.length {
			progressInSong = float64(sampleInSong) / float64(currentSong.length)
		}
		songLength = time.Duration(currentSong.length) * time.Second / time.Duration(schedule.SampleRate) / time.Nanosecond
	}

	return &playbackInformation{
		CurrentSong:         currentSong,
		CurrentSample:       sample,
		PausesInCurrentSong: pausesInCurrentSong,
		Now:                 now,
		Playing:             playing,
		Volume:              s.Volume,
		SongLength:          songLength,
		TimeInSong:          timeInSong,
		ProgressInSong:      progressInSong}
}

type pausesInCurrentSongState struct {
	playing             bool
	pausesInCurrentSong int64
	pauseBegin          int64
}

func (picss *pausesInCurrentSongState) handlePause(sample int64, p pauseToggle, startIndex uint64) {
	if uint64(sample) <= p.toggleIndex || picss.playing == p.playing {
		return
	}

	if p.playing {
		picss.pausesInCurrentSong += int64(p.toggleIndex) - picss.pauseBegin
	} else {
		picss.pauseBegin = int64(p.toggleIndex)
	}

	if p.toggleIndex < startIndex {
		if p.playing {
			picss.pausesInCurrentSong = 0
		} else {
			picss.pausesInCurrentSong = int64(p.toggleIndex) - int64(startIndex)
		}
	}

	picss.playing = p.playing
}

func (s *state) pausesInCurrentSong(sample int64, currentSong upcomingSong) (pausesInCurrentSong int64, playing bool) {
	s.PausesMutex.RLock()
	defer s.PausesMutex.RUnlock()

	picss := &pausesInCurrentSongState{playing: true, pausesInCurrentSong: 0, pauseBegin: 0}
	for _, p := range s.Pauses {
		picss.handlePause(sample, p, currentSong.startIndex)
	}

	if !picss.playing {
		picss.pausesInCurrentSong += sample - picss.pauseBegin
	}

	return picss.pausesInCurrentSong, picss.playing
}

func (s *state) removeOldPauses(currentSong upcomingSong) {
	s.PausesMutex.Lock()
	defer s.PausesMutex.Unlock()

	passed := 0
	for i, p := range s.Pauses {
		if p.playing && p.toggleIndex < currentSong.startIndex {
			passed = i
		}
	}
	if 0 < passed {
		copy(s.Pauses, s.Pauses[passed:])
		for i := len(s.Pauses) - passed; i < len(s.Pauses); i++ {
			s.Pauses[i] = pauseToggle{}
		}
		s.Pauses = s.Pauses[:len(s.Pauses)-passed]
	}
}

func (s *state) currentSong(sample int64) (currentSong upcomingSong) {
	s.SongsMutex.Lock()
	defer s.SongsMutex.Unlock()

	currentSong = upcomingSong{filename: "None", startIndex: 0, length: 0}
	for i := len(s.Songs) - 1; 0 <= i; i-- {
		if int64(s.Songs[i].startIndex) < sample {
			currentSong = s.Songs[i]
			break
		}
	}

	return
}

func (s *state) currentSample(now int64) (sample int64) {
	sample = 0
	s.removeOldChunks(now)

	s.ChunksMutex.Lock()
	defer s.ChunksMutex.Unlock()

	if len(s.Chunks) != 0 {
		timeInChunk := now - s.Chunks[0].startTime
		sampleInChunk := int64(time.Duration(timeInChunk) * time.Nanosecond * time.Duration(schedule.SampleRate) / time.Second)
		sample = int64(s.Chunks[0].startIndex) + sampleInChunk
	}

	return
}

func (s *state) removeOldChunks(now int64) {
	s.ChunksMutex.Lock()
	defer s.ChunksMutex.Unlock()

	if len(s.Chunks) == 0 {
		return
	}

	passed := 0
	for ; passed < len(s.Chunks); passed++ {
		if now <= s.Chunks[passed].endTime() {
			break
		}
	}
	if passed != 0 {
		copy(s.Chunks, s.Chunks[passed:])
		for i := len(s.Chunks) - passed; i < len(s.Chunks); i++ {
			s.Chunks[i] = upcomingChunk{}
		}
		s.Chunks = s.Chunks[:len(s.Chunks)-passed]
	}
}

type playbackInformation struct {
	CurrentSong         upcomingSong
	CurrentSample       int64
	PausesInCurrentSong int64
	Now                 int64
	Playing             bool
	Volume              float64
	SongLength          time.Duration
	TimeInSong          time.Duration
	ProgressInSong      float64
}

func (pbi playbackInformation) playingString() string {
	if pbi.Playing {
		return "Playing"
	}
	return "Paused"
}

var currentState = &state{Songs: make([]upcomingSong, 0), Chunks: make([]upcomingChunk, 0)}

type upcomingSong struct {
	filename   string
	startIndex uint64
	length     int64
	lyrics     []metadata.LyricsLine
	metadata   metadata.SongMetadata
}

type upcomingChunk struct {
	startTime  int64
	startIndex uint64
	size       uint64
}

func (uc upcomingChunk) endTime() int64 {
	return uc.startTime + uc.length()
}

func (uc upcomingChunk) length() int64 {
	return int64(time.Duration(uc.size)*time.Second) / int64(schedule.SampleRate)
}

type pauseToggle struct {
	playing     bool
	toggleIndex uint64
}
