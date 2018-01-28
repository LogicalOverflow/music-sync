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

var TimeSyncInterval = 10 * time.Minute
var TimeSyncCycles = 500
var TimeSyncCycleDelay = 10 * time.Millisecond
var StreamChunkSize = 44100 * 4
var StreamChunkTime = 4 * time.Second
var NanBreakSize = 44100 * 1
var StreamStartDelay = 5 * time.Second
var StreamDelay = 15 * time.Second
var SampleRate = 44100

func Server(sender comm.MessageSender) {
	playlist := playback.NewPlaylist(SampleRate, []string{}, NanBreakSize)
	volume := 0.1
	comm.NewSlaveHandler = func(c comm.Channel, s comm.MessageSender) {
		switch c {
		case comm.Channel_AUDIO:
			s.SendMessage(&comm.SetVolumeRequest{Volume: volume})
		}
	}

	go playlist.StreamLoop()

	go func() {
		time.Sleep(StreamStartDelay)
		start := timing.GetSyncedTime() + int64(StreamDelay/time.Nanosecond)
		index := int64(0)
		for range time.Tick(StreamChunkTime) {
			low := make([]float64, StreamChunkSize)
			high := make([]float64, StreamChunkSize)

			playlist.Fill(low, high)

			sender.SendMessage(&comm.QueueChunkRequest{
				StartTime:  start + int64(index)*int64(StreamChunkTime/time.Nanosecond),
				ChunkId:    index,
				SampleLow:  low,
				SampleHigh: high,
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
				if pos, err := strconv.Atoi(args[1]); err != nil {
					return "", false
				} else {
					playlist.InsertSong(song, pos)
				}
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
			} else {
				return fmt.Sprintf("setting volume to %.3f", volume), true
			}
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

func Player(sender comm.MessageSender) {
	go func() {
		if err := playback.Init(SampleRate); err != nil {
			logger.Fatalf("failed to initialized playback: %v", err)
			os.Exit(1)
		}
	}()

	go func() {
		SyncTime(sender)
		for range time.Tick(TimeSyncInterval) {
			SyncTime(sender)
		}
	}()

	go func() {
		if err := sender.SendMessage(&comm.SubscribeChannelRequest{Channel: comm.Channel_AUDIO}); err != nil {
			logger.Errorf("failed to subscribe to audio channel")
			os.Exit(1)
		}
	}()
}

func SyncTime(sender comm.MessageSender) {
	logger.Debugf("syncing time")
	timing.ResetOffsets(TimeSyncCycles)
	for i := 0; i < TimeSyncCycles; i++ {
		if err := sender.SendMessage(&comm.TimeSyncRequest{ClientSend: timing.GetRawTime()}); err != nil {
			logger.Warnf("failed to send sync time request: %v", err)
		}
		time.Sleep(TimeSyncCycleDelay)
	}
}
