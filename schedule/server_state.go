package schedule

import (
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/metadata"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/timing"
	"sync"
	"time"
)

type serverState struct {
	sender comm.MessageSender

	lyricsProvider   metadata.LyricsProvider
	metadataProvider metadata.Provider

	playlist *playback.Playlist
	volume   float64

	newestSong *comm.NewSongInfo

	pauses      []*comm.PauseInfo
	pausesMutex sync.RWMutex
}

func (ss *serverState) sendVolume(s comm.MessageSender) {
	s.SendMessage(&comm.SetVolumeRequest{Volume: ss.volume})
}

func (ss *serverState) sendNewestSong(s comm.MessageSender) {
	if ss.newestSong != nil {
		s.SendMessage(ss.newestSong)
	}
}

func (ss *serverState) sendPauses(s comm.MessageSender) {
	ss.pausesMutex.RLock()
	for _, p := range ss.pauses {
		s.SendMessage(p)
	}
	ss.pausesMutex.RUnlock()
}

func (ss *serverState) createClientHandler() func(comm.Channel, comm.MessageSender) {
	return func(c comm.Channel, s comm.MessageSender) {
		switch c {
		case comm.Channel_AUDIO:
			ss.sendVolume(s)
		case comm.Channel_META:
			ss.sendVolume(s)
			ss.sendNewestSong(s)
			ss.sendPauses(s)
		}
	}
}

func toWireLyrics(lyrics []metadata.LyricsLine) []*comm.NewSongInfo_SongLyricsLine {
	wireLyrics := make([]*comm.NewSongInfo_SongLyricsLine, len(lyrics))
	for i, l := range lyrics {
		wireLine := make([]*comm.NewSongInfo_SongLyricsAtom, len(l))
		for j, a := range l {
			wireLine[j] = &comm.NewSongInfo_SongLyricsAtom{
				Timestamp: a.Timestamp,
				Caption:   a.Caption,
			}
		}
		wireLyrics[i] = &comm.NewSongInfo_SongLyricsLine{Atoms: wireLine}
	}
	return wireLyrics
}

func (ss *serverState) createNewSongHandler() func(uint64, string, int64) {
	return func(startSampleIndex uint64, filename string, songLength int64) {
		lyrics := ss.lyricsProvider.CollectLyrics(filename)
		wireLyrics := toWireLyrics(lyrics)

		md := ss.metadataProvider.CollectMetadata(filename)

		ss.newestSong = &comm.NewSongInfo{
			FirstSampleOfSongIndex: startSampleIndex,
			SongFileName:           filename,
			SongLength:             songLength,
			Lyrics:                 wireLyrics,
			Metadata: &comm.NewSongInfo_SongMetadata{
				Title:  md.Title,
				Artist: md.Artist,
				Album:  md.Album,
			},
		}
		ss.sender.SendMessage(ss.newestSong)
	}
}

func (ss *serverState) createPauseToggleHandler() func(bool, uint64) {
	return func(playing bool, sample uint64) {
		pause := &comm.PauseInfo{
			Playing:           playing,
			ToggleSampleIndex: sample,
		}
		go func(pause *comm.PauseInfo) {
			ss.pausesMutex.Lock()
			defer ss.pausesMutex.Unlock()
			ss.pauses = append(ss.pauses, pause)
			go ss.removeOldPauses()
		}(pause)
		ss.sender.SendMessage(pause)
	}
}

func (ss *serverState) removablePauses() int {
	if ss.newestSong == nil {
		return 0
	}

	ss.pausesMutex.Lock()
	defer ss.pausesMutex.Unlock()
	passed := 0
	for i, p := range ss.pauses {
		if p.ToggleSampleIndex < ss.newestSong.FirstSampleOfSongIndex && p.Playing {
			passed = i
		} else if ss.newestSong.FirstSampleOfSongIndex < p.ToggleSampleIndex {
			break
		}
	}
	return passed
}

func (ss *serverState) removeOldPauses() {
	removable := ss.removablePauses()
	if 0 < removable {
		ss.pausesMutex.Lock()
		defer ss.pausesMutex.Unlock()

		copy(ss.pauses, ss.pauses[removable:])
		for i := len(ss.pauses) - removable; i < len(ss.pauses); i++ {
			ss.pauses[i] = nil
		}
		ss.pauses = ss.pauses[:len(ss.pauses)-removable]
	}
}

func (ss *serverState) streamMusic() {
	time.Sleep(StreamStartDelay)
	start := timing.GetSyncedTime() + int64(StreamDelay/time.Nanosecond)
	index := int64(0)
	for range time.Tick(StreamChunkTime) {
		low := make([]float64, StreamChunkSize)
		high := make([]float64, StreamChunkSize)

		firstSampleIndex := ss.playlist.Fill(low, high)

		go ss.sender.SendMessage(&comm.QueueChunkRequest{
			StartTime:        start + int64(index)*int64(StreamChunkTime/time.Nanosecond),
			ChunkId:          index,
			SampleLow:        low,
			SampleHigh:       high,
			FirstSampleIndex: firstSampleIndex,
		})
		go ss.sender.SendMessage(&comm.ChunkInfo{
			StartTime:        start + int64(index)*int64(StreamChunkTime/time.Nanosecond),
			FirstSampleIndex: firstSampleIndex,
			ChunkSize:        uint64(StreamChunkSize),
		})
		index++
	}
}
