package schedule

import (
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/metadata"
)

var testSongMessages = []*comm.NewSongInfo{
	{FirstSampleOfSongIndex: 0, SongFileName: "empty-song", SongLength: 0},
	{FirstSampleOfSongIndex: 16, SongFileName: "a-song", SongLength: 64},
}

var testPauseMessages = [][]*comm.PauseInfo{
	{}, {{Playing: true, ToggleSampleIndex: 1}}, {{Playing: false, ToggleSampleIndex: 2}},
	{{Playing: false, ToggleSampleIndex: 3}, {Playing: true, ToggleSampleIndex: 4}},
}

var removePauseCases = []struct {
	songStartSample uint64
	pauses          []*comm.PauseInfo
	removable       int
	result          []*comm.PauseInfo
}{
	{
		songStartSample: 0,
		pauses:          []*comm.PauseInfo{},
		removable:       0,
		result:          []*comm.PauseInfo{},
	},
	{
		songStartSample: 0,
		pauses:          []*comm.PauseInfo{{Playing: true, ToggleSampleIndex: 8}, {Playing: false, ToggleSampleIndex: 16}},
		removable:       0,
		result:          []*comm.PauseInfo{{Playing: true, ToggleSampleIndex: 8}, {Playing: false, ToggleSampleIndex: 16}},
	},
	{
		songStartSample: 0,
		pauses:          []*comm.PauseInfo{{Playing: false, ToggleSampleIndex: 8}, {Playing: true, ToggleSampleIndex: 16}},
		removable:       0,
		result:          []*comm.PauseInfo{{Playing: false, ToggleSampleIndex: 8}, {Playing: true, ToggleSampleIndex: 16}},
	},
	{
		songStartSample: 0,
		pauses:          []*comm.PauseInfo{{Playing: false, ToggleSampleIndex: 8}, {Playing: true, ToggleSampleIndex: 16}, {Playing: false, ToggleSampleIndex: 24}},
		removable:       0,
		result:          []*comm.PauseInfo{{Playing: false, ToggleSampleIndex: 8}, {Playing: true, ToggleSampleIndex: 16}, {Playing: false, ToggleSampleIndex: 24}},
	},
	{
		songStartSample: 128,
		pauses:          []*comm.PauseInfo{{Playing: true, ToggleSampleIndex: 8}, {Playing: false, ToggleSampleIndex: 16}},
		removable:       0,
		result:          []*comm.PauseInfo{{Playing: true, ToggleSampleIndex: 8}, {Playing: false, ToggleSampleIndex: 16}},
	},
	{
		songStartSample: 128,
		pauses:          []*comm.PauseInfo{{Playing: false, ToggleSampleIndex: 8}, {Playing: true, ToggleSampleIndex: 16}},
		removable:       1,
		result:          []*comm.PauseInfo{{Playing: true, ToggleSampleIndex: 16}},
	},
	{
		songStartSample: 128,
		pauses:          []*comm.PauseInfo{{Playing: false, ToggleSampleIndex: 8}, {Playing: true, ToggleSampleIndex: 16}, {Playing: false, ToggleSampleIndex: 24}},
		removable:       1,
		result:          []*comm.PauseInfo{{Playing: true, ToggleSampleIndex: 16}, {Playing: false, ToggleSampleIndex: 24}},
	},
}

var toWireLyricsCases = []struct {
	lyrics []metadata.LyricsLine
	result []*comm.NewSongInfo_SongLyricsLine
}{
	{
		lyrics: []metadata.LyricsLine{},
		result: []*comm.NewSongInfo_SongLyricsLine{},
	},
	{
		lyrics: []metadata.LyricsLine{
			{metadata.LyricsAtom{Timestamp: 1, Caption: "caption-1-1"}, metadata.LyricsAtom{Timestamp: 2, Caption: "caption-1-2"}, metadata.LyricsAtom{Timestamp: 3, Caption: "caption-1-3"}},
			{metadata.LyricsAtom{Timestamp: 4, Caption: "caption-2-1"}, metadata.LyricsAtom{Timestamp: 5, Caption: "caption-2-2"}, metadata.LyricsAtom{Timestamp: 6, Caption: "caption-2-3"}},
			{metadata.LyricsAtom{Timestamp: 7, Caption: "caption-3-1"}, metadata.LyricsAtom{Timestamp: 8, Caption: "caption-3-2"}, metadata.LyricsAtom{Timestamp: 9, Caption: "caption-3-3"}},
		},
		result: []*comm.NewSongInfo_SongLyricsLine{
			{Atoms: []*comm.NewSongInfo_SongLyricsAtom{{Timestamp: 1, Caption: "caption-1-1"}, {Timestamp: 2, Caption: "caption-1-2"}, {Timestamp: 3, Caption: "caption-1-3"}}},
			{Atoms: []*comm.NewSongInfo_SongLyricsAtom{{Timestamp: 4, Caption: "caption-2-1"}, {Timestamp: 5, Caption: "caption-2-2"}, {Timestamp: 6, Caption: "caption-2-3"}}},
			{Atoms: []*comm.NewSongInfo_SongLyricsAtom{{Timestamp: 7, Caption: "caption-3-1"}, {Timestamp: 8, Caption: "caption-3-2"}, {Timestamp: 9, Caption: "caption-3-3"}}},
		},
	},
}
