// Package schedule contains methods to start different types of clients/servers
package schedule

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/ssh"
	"github.com/LogicalOverflow/music-sync/timing"
	"github.com/LogicalOverflow/music-sync/util"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger = log.GetLogger("shed")

// TimeSyncInterval is the time interval between syncing time to server
var TimeSyncInterval = 10 * time.Minute
// TimeSyncCycles is the number of cycles used to sync time to server
var TimeSyncCycles = 500
// TimeSyncCycleDelay is the delay between cycles in one time sync
var TimeSyncCycleDelay = 10 * time.Millisecond
// StreamChunkSize is the size of one stream chunk in samples
var StreamChunkSize = 44100 * 4
// StreamChunkTime is the duration is takes to play one stream chunk
var StreamChunkTime = 4 * time.Second
// NanBreakSize is the number of nan-samples to insert between songs, which players use to realign playback
var NanBreakSize = 44100 * 1
// StreamStartDelay is the delay before starting the stream
var StreamStartDelay = 5 * time.Second
// StreamDelay is the delay of the stream, which players use to decode chunks
var StreamDelay = 15 * time.Second
// SampleRate is the sample rate of the stream
var SampleRate = 44100

// Server starts a music-sync server, using sender to communicate with all clients
func Server(sender comm.MessageSender) {
	playlist := playback.NewPlaylist(SampleRate, []string{}, NanBreakSize)
	volume := 0.1
	comm.NewClientHandler = func(c comm.Channel, s comm.MessageSender) {
		switch c {
		case comm.Channel_AUDIO:
			s.SendMessage(&comm.SetVolumeRequest{Volume: volume})
		}
	}

	go playlist.StreamLoop()

	playlist.SetNewSongHandler(func(startSampleIndex uint64, filename string) {
		sender.SendMessage(&comm.NewSongInfo{
			FirstSampleOfSongIndex: startSampleIndex,
			SongFileName:           filename,
		})
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
			song := args[0]
			if len(args) == 1 {
				playlist.AddSong(song)
			} else {
				pos, err := strconv.Atoi(args[1])
				if err != nil {
					return "", false
				}
				playlist.InsertSong(song, pos)
			}
			return fmt.Sprintf("song %s added to playlist", song), true
		},
		Options: func(prefix string, arg int) []string {
			if arg != 0 {
				return []string{}
			}
			songs := util.ListAllSongs(playback.AudioDir, "")
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

// Player starts a music-sync player, using sender to communicate with the server
func Player(sender comm.MessageSender) {
	go func() {
		if err := playback.Init(SampleRate); err != nil {
			logger.Fatalf("failed to initialized playback: %v", err)
			os.Exit(1)
		}
	}()

	go func() {
		syncTime(sender)
		for range time.Tick(TimeSyncInterval) {
			syncTime(sender)
		}
	}()

	go func() {
		if err := sender.SendMessage(&comm.SubscribeChannelRequest{Channel: comm.Channel_AUDIO}); err != nil {
			logger.Errorf("failed to subscribe to audio channel")
			os.Exit(1)
		}
	}()
}

// Infoer start a music-sync client in infoer mode, using sender to communicate with the server
func Infoer(sender comm.MessageSender) {
	go func() {
		syncTime(sender)
		for range time.Tick(TimeSyncInterval) {
			syncTime(sender)
		}
	}()

	go func() {
		if err := sender.SendMessage(&comm.SubscribeChannelRequest{Channel: comm.Channel_META}); err != nil {
			logger.Errorf("failed to subscribe to meta channel")
			os.Exit(1)
		}
	}()
}

func syncTime(sender comm.MessageSender) {
	logger.Infof("syncing time")
	timing.ResetOffsets(TimeSyncCycles)
	for i := 0; i < TimeSyncCycles; i++ {
		if err := sender.SendMessage(&comm.TimeSyncRequest{ClientSend: timing.GetRawTime()}); err != nil {
			logger.Warnf("failed to send sync time request: %v", err)
		}
		time.Sleep(TimeSyncCycleDelay)
	}
}
