package schedule

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/metadata"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/ssh"
	"github.com/LogicalOverflow/music-sync/timing"
	"github.com/LogicalOverflow/music-sync/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Server starts a music-sync server, using sender to communicate with all clients
func Server(sender comm.MessageSender) {
	lyricsProvider := metadata.GetLyricsProvider()
	metadataProvider := metadata.GetProvider()

	playlist := playback.NewPlaylist(SampleRate, []string{}, NanBreakSize)
	volume := 0.1
	var newestSong *comm.NewSongInfo

	pauses := make([]*comm.PauseInfo, 0)
	var pausesMutex sync.RWMutex

	comm.NewClientHandler = func(c comm.Channel, s comm.MessageSender) {
		switch c {
		case comm.Channel_AUDIO:
			s.SendMessage(&comm.SetVolumeRequest{Volume: volume})
		case comm.Channel_META:
			s.SendMessage(&comm.SetVolumeRequest{Volume: volume})
			if newestSong != nil {
				sender.SendMessage(newestSong)
			}
			pausesMutex.RLock()
			for _, p := range pauses {
				sender.SendMessage(p)
			}
			pausesMutex.RUnlock()
		}
	}

	go playlist.StreamLoop()

	playlist.SetNewSongHandler(func(startSampleIndex uint64, filename string, songLength int64) {
		lyrics := lyricsProvider.CollectLyrics(filename)
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

		md := metadataProvider.CollectMetadata(filename)

		newestSong = &comm.NewSongInfo{
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
		sender.SendMessage(newestSong)
	})
	playlist.SetPauseToggleHandler(func(playing bool, sample uint64) {
		pause := &comm.PauseInfo{
			Playing:           playing,
			ToggleSampleIndex: sample,
		}
		go func(pause *comm.PauseInfo) {
			pausesMutex.Lock()
			pauses = append(pauses, pause)
			if newestSong != nil {
				passed := 0
				for i, p := range pauses {
					if p.ToggleSampleIndex < newestSong.FirstSampleOfSongIndex && p.Playing {
						passed = i
					} else if newestSong.FirstSampleOfSongIndex < p.ToggleSampleIndex {
						break
					}
				}
				if 0 < passed {
					copy(pauses, pauses[passed:])
					for i := len(pauses) - passed; i < len(pauses); i++ {
						pauses[i] = nil
					}
					pauses = pauses[:len(pauses)-passed]
				}
			}
			pausesMutex.Unlock()
		}(pause)
		sender.SendMessage(pause)
	})

	go func() {
		time.Sleep(StreamStartDelay)
		start := timing.GetSyncedTime() + int64(StreamDelay/time.Nanosecond)
		index := int64(0)
		for range time.Tick(StreamChunkTime) {
			low := make([]float64, StreamChunkSize)
			high := make([]float64, StreamChunkSize)

			firstSampleIndex := playlist.Fill(low, high)

			go sender.SendMessage(&comm.QueueChunkRequest{
				StartTime:        start + int64(index)*int64(StreamChunkTime/time.Nanosecond),
				ChunkId:          index,
				SampleLow:        low,
				SampleHigh:       high,
				FirstSampleIndex: firstSampleIndex,
			})
			go sender.SendMessage(&comm.ChunkInfo{
				StartTime:        start + int64(index)*int64(StreamChunkTime/time.Nanosecond),
				FirstSampleIndex: firstSampleIndex,
				ChunkSize:        uint64(StreamChunkSize),
			})
			index++
		}
	}()

	ssh.RegisterCommand(ssh.Command{
		Name:  "queue",
		Usage: "filename [position in playlist]",
		Info:  "adds a song to the playlist",
		Exec: func(args []string) (string, bool) {
			if len(args) != 2 && len(args) != 1 {
				return "", false
			}
			songPattern := args[0]
			songs, err := util.ListGlobFiles(playback.AudioDir, songPattern)
			if err != nil {
				return fmt.Sprintf("glob pattern is invalid: %v", err), true
			}
			songs = util.FilterSongs(songs)
			if len(songs) == 0 {
				return fmt.Sprintf("no song matches the glob pattern %s", songPattern), true
			}

			if len(args) == 1 {
				for _, s := range songs {
					playlist.AddSong(s)
				}
			} else {
				pos, err := strconv.Atoi(args[1])
				if err != nil {
					return "", false
				}
				for i, s := range songs {
					playlist.InsertSong(s, pos+i)
				}
			}
			return fmt.Sprintf("%d song(s) added to playlist: %s", len(songs), strings.Join(songs, ", ")), true
		},
		Options: func(prefix string, arg int) []string {
			if arg != 0 {
				return []string{}
			}
			songs := util.FilterSongs(util.ListAllFiles(playback.AudioDir, ""))
			options := make([]string, 0, len(songs))
			for _, song := range songs {
				if strings.HasPrefix(song, prefix) {
					options = append(options, song)
				}
			}
			return options
		},
	})
	ssh.RegisterCommand(ssh.Command{
		Name:  "playlist",
		Usage: "",
		Info:  "prints the current playlist",
		Exec: func([]string) (string, bool) {
			songs := playlist.Songs()
			entries := make([]string, len(songs))
			format := fmt.Sprintf("  [%%0%dd] %%s", len(strconv.Itoa(len(songs)-1)))
			for i, s := range songs {
				entries[i] = fmt.Sprintf(format, i, s)
			}
			var playingStatus string
			if playlist.Playing() {
				playingStatus = "Playing"
			} else {
				playingStatus = "Paused"
			}

			songList := "Empty"
			if 0 < len(entries) {
				songList = "\n" + strings.Join(entries, "\n")
			}
			currSong := playlist.CurrentSong()
			if currSong == "" {
				currSong = "None"
			}

			return fmt.Sprintf("Current Playlist (%s): %s\nCurrent Song: %s", playingStatus, songList, currSong), true
		},
	})
	ssh.RegisterCommand(ssh.Command{
		Name:  "remove",
		Usage: "position",
		Info:  "removes a song from the playlist",
		Exec: func(args []string) (string, bool) {
			if len(args) != 1 {
				return "", false
			}
			pos, err := strconv.Atoi(args[0])
			if err != nil {
				return "", false
			}
			song := playlist.RemoveSong(pos)
			return fmt.Sprintf("removed song %s at position %d from playlist", song, pos), true
		},
	})
	ssh.RegisterCommand(ssh.Command{
		Name:  "jump",
		Usage: "position",
		Info:  "jumps in the playlist",
		Exec: func(args []string) (string, bool) {
			if len(args) != 1 {
				return "", false
			}
			pos, err := strconv.Atoi(args[0])
			if err != nil {
				return "", false
			}
			playlist.SetPos(pos)
			return "jumped", true
		},
	})
	ssh.RegisterCommand(ssh.Command{
		Name:  "volume",
		Usage: "volume",
		Info:  "set the playback volume",
		Exec: func(args []string) (string, bool) {
			if len(args) != 1 {
				return "", false
			}
			var err error
			volume, err = strconv.ParseFloat(args[0], 64)
			if err != nil {
				return "", false
			}
			if err := sender.SendMessage(&comm.SetVolumeRequest{Volume: volume}); err != nil {
				return fmt.Sprintf("failed to set volume to %.3f: %v", volume, err), true
			}
			return fmt.Sprintf("setting volume to %.3f", volume), true
		},
	})
	ssh.RegisterCommand(ssh.Command{
		Name:  "pause",
		Usage: "",
		Info:  "pauses playback",
		Exec: func([]string) (string, bool) {
			playlist.SetPlaying(false)
			return "playback paused", true
		},
	})
	ssh.RegisterCommand(ssh.Command{
		Name:  "resume",
		Usage: "",
		Info:  "resumes playback",
		Exec: func([]string) (string, bool) {
			playlist.SetPlaying(true)
			return "playback resumed", true
		},
	})
}
