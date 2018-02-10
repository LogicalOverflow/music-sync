package main

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/metadata"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testPackageLyrics = []*comm.NewSongInfo_SongLyricsLine{
	{Atoms: []*comm.NewSongInfo_SongLyricsAtom{{Timestamp: 1, Caption: "caption-0-0"}, {Timestamp: 2, Caption: "caption-0-1"}}},
	{Atoms: []*comm.NewSongInfo_SongLyricsAtom{{Timestamp: 3, Caption: "caption-1-0"}, {Timestamp: 4, Caption: "caption-1-1"}}},
}
var testPackageMetadata = &comm.NewSongInfo_SongMetadata{
	Title:  "song-title",
	Artist: "song-artist",
	Album:  "song-album",
}

var testMetadataLyrics = []metadata.LyricsLine{
	{{Timestamp: 1, Caption: "caption-0-0"}, {Timestamp: 2, Caption: "caption-0-1"}},
	{{Timestamp: 3, Caption: "caption-1-0"}, {Timestamp: 4, Caption: "caption-1-1"}},
}
var testMetadata = metadata.SongMetadata{
	Title:  "song-title",
	Artist: "song-artist",
	Album:  "song-album",
}

func TestInfoerPackageHandler_HandleChunkInfo(t *testing.T) {
	ph := newInfoerPackageHandler()
	infos := make([]upcomingChunk, 0)
	currentState.Chunks = make([]upcomingChunk, 0)
	for i := 15; 0 <= i; i-- {
		info := &comm.ChunkInfo{
			StartTime:        int64(i) * 1e9,
			FirstSampleIndex: uint64(i) * 512,
			ChunkSize:        512,
		}
		ph.HandleChunkInfo(info, nil)

		infos = append([]upcomingChunk{{startTime: info.StartTime, startIndex: info.FirstSampleIndex, size: info.ChunkSize}}, infos...)
		assert.Equal(t, infos, currentState.Chunks, "HandleChunkInfo did not add to currentState Chunks correctly")
	}
}

func TestInfoerPackageHandler_HandleNewSongInfo(t *testing.T) {
	ph := newInfoerPackageHandler()
	songs := make([]upcomingSong, 0)
	currentState.Songs = make([]upcomingSong, 0)
	for i := 15; 0 <= i; i-- {
		song := &comm.NewSongInfo{
			FirstSampleOfSongIndex: uint64(i) * 512,
			SongFileName:           fmt.Sprintf("song-%02d", i),
			SongLength:             256,
			Lyrics:                 testPackageLyrics,
			Metadata:               testPackageMetadata,
		}
		ph.HandleNewSongInfo(song, nil)

		songs = append([]upcomingSong{{filename: song.SongFileName, startIndex: song.FirstSampleOfSongIndex, length: song.SongLength, lyrics: testMetadataLyrics, metadata: testMetadata}}, songs...)
		assert.Equal(t, songs, currentState.Songs, "HandleNewSongInfo did not add to currentState Songs correctly")
	}
}

func TestInfoerPackageHandler_HandlePauseInfo(t *testing.T) {
	ph := newInfoerPackageHandler()
	pauses := make([]pauseToggle, 0)
	currentState.Pauses = make([]pauseToggle, 0)
	for i := 15; 0 <= i; i-- {
		pause := &comm.PauseInfo{
			Playing:           i%2 == 0,
			ToggleSampleIndex: uint64(i) * 512,
		}
		ph.HandlePauseInfo(pause, nil)

		pauses = append([]pauseToggle{{playing: pause.Playing, toggleIndex: pause.ToggleSampleIndex}}, pauses...)
		assert.Equal(t, pauses, currentState.Pauses, "HandlePauseInfo did not add to currentState Pauses correctly")
	}
}

func TestInfoerPackageHandler_HandleSetVolumeRequest(t *testing.T) {
	ph := newInfoerPackageHandler()
	currentState.Volume = -1
	for i := float64(0); i <= 1; i += 0.125 {
		ph.HandleSetVolumeRequest(&comm.SetVolumeRequest{Volume: i}, nil)
		assert.Equal(t, i, currentState.Volume, "HandleSetVolumeRequest did not update currentState volume correctly")
	}
}
