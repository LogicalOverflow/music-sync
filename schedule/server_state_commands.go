package schedule

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/ssh"
	"github.com/LogicalOverflow/music-sync/util"
	"strconv"
	"strings"
)

func (ss *serverState) queueCommand() ssh.Command {
	return ssh.Command{
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
					ss.playlist.AddSong(s)
				}
			} else {
				pos, err := strconv.Atoi(args[1])
				if err != nil {
					return "", false
				}
				for i, s := range songs {
					ss.playlist.InsertSong(s, pos+i)
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
	}
}

func (ss *serverState) playlistCommand() ssh.Command {
	return ssh.Command{
		Name:  "playlist",
		Usage: "",
		Info:  "prints the current playlist",
		Exec: func([]string) (string, bool) {
			songs := ss.playlist.Songs()
			entries := make([]string, len(songs))
			format := fmt.Sprintf("  [%%0%dd] %%s", len(strconv.Itoa(len(songs)-1)))
			for i, s := range songs {
				entries[i] = fmt.Sprintf(format, i, s)
			}
			var playingStatus string
			if ss.playlist.Playing() {
				playingStatus = "Playing"
			} else {
				playingStatus = "Paused"
			}

			songList := "Empty"
			if 0 < len(entries) {
				songList = "\n" + strings.Join(entries, "\n")
			}
			currSong := ss.playlist.CurrentSong()
			if currSong == "" {
				currSong = "None"
			}

			return fmt.Sprintf("Current Playlist (%s): %s\nCurrent Song: %s", playingStatus, songList, currSong), true
		},
	}
}

func (ss *serverState) removeCommand() ssh.Command {
	return ssh.Command{
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
			song := ss.playlist.RemoveSong(pos)
			return fmt.Sprintf("removed song %s at position %d from playlist", song, pos), true
		},
	}
}

func (ss *serverState) jumpCommand() ssh.Command {
	return ssh.Command{
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
			ss.playlist.SetPos(pos)
			return "jumped", true
		},
	}
}

func (ss *serverState) volumeCommand() ssh.Command {
	return ssh.Command{
		Name:  "volume",
		Usage: "volume",
		Info:  "set the playback volume",
		Exec: func(args []string) (string, bool) {
			if len(args) != 1 {
				return "", false
			}
			var err error
			ss.volume, err = strconv.ParseFloat(args[0], 64)
			if err != nil {
				return "", false
			}
			if err := ss.sender.SendMessage(&comm.SetVolumeRequest{Volume: ss.volume}); err != nil {
				return fmt.Sprintf("failed to set volume to %.3f: %v", ss.volume, err), true
			}
			return fmt.Sprintf("setting volume to %.3f", ss.volume), true
		},
	}
}

func (ss *serverState) pauseCommand() ssh.Command {
	return ssh.Command{
		Name:  "pause",
		Usage: "",
		Info:  "pauses playback",
		Exec: func([]string) (string, bool) {
			ss.playlist.SetPlaying(false)
			return "playback paused", true
		},
	}
}

func (ss *serverState) resumeCommand() ssh.Command {
	return ssh.Command{
		Name:  "resume",
		Usage: "",
		Info:  "resumes playback",
		Exec: func([]string) (string, bool) {
			ss.playlist.SetPlaying(true)
			return "playback resumed", true
		},
	}
}
