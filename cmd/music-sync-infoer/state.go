package main

import (
	"github.com/LogicalOverflow/music-sync/metadata"
	"github.com/LogicalOverflow/music-sync/schedule"
	"github.com/LogicalOverflow/music-sync/timing"
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

func (s *state) Info() *playbackInformation {
	now := timing.GetSyncedTime()
	sample := s.currentSample(now)

	currentSong := s.currentSong(sample)

	s.removeOldPauses(currentSong)

	pausesInCurrentSong, playing := s.pausesInCurrentSong(sample, currentSong)

	return &playbackInformation{
		CurrentSong:         currentSong,
		CurrentSample:       sample,
		PausesInCurrentSong: pausesInCurrentSong,
		Now:                 now,
		Playing:             playing,
		Volume:              s.Volume,
	}
}

func (s *state) pausesInCurrentSong(sample int64, currentSong upcomingSong) (pausesInCurrentSong int64, playing bool) {
	s.PausesMutex.RLock()
	defer s.PausesMutex.RUnlock()

	playing = true
	pausesInCurrentSong = int64(0)

	pauseBegin := int64(0)
	for _, p := range s.Pauses {
		if p.toggleIndex < uint64(sample) && playing != p.playing {
			if p.playing {
				pausesInCurrentSong += int64(p.toggleIndex) - pauseBegin
			} else {
				pauseBegin = int64(p.toggleIndex)
			}

			if p.toggleIndex < currentSong.startIndex {
				if p.playing {
					pausesInCurrentSong = 0
				} else {
					pausesInCurrentSong = int64(p.toggleIndex) - int64(currentSong.startIndex)
				}
			}

			playing = p.playing
		}
	}

	if !playing {
		pausesInCurrentSong += sample - pauseBegin
	}

	return
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
	sample = -1
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
	for ; s.Chunks[passed].endTime() < now; passed++ {
		if len(s.Chunks) <= passed {
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
