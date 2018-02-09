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

func parseFloatParam(args []string, index int) (float64, bool) {
	if len(args) <= index {
		return 0, false
	}
	v, err := strconv.ParseFloat(args[index], 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func parseStringParam(args []string, index int) (string, bool) {
	if len(args) <= index {
		return "", false
	}
	return args[index], true
}

func parseIntParam(args []string, index int) (int, bool) {
	if len(args) <= index {
		return 0, false
	}
	v, err := strconv.Atoi(args[index])
	if err != nil {
		return 0, false
	}
	return v, true
}

func (ss *serverState) queueCommandExec(args []string) (string, bool) {
	songPattern, ok := parseStringParam(args, 0)
	if !ok {
		return "", false
	}

	songs, err := util.ListGlobFiles(playback.AudioDir, songPattern)
	if err != nil {
		return fmt.Sprintf("glob pattern is invalid: %v", err), true
	}
	songs = util.FilterSongs(songs)
	if len(songs) == 0 {
		return fmt.Sprintf("no song matches the glob pattern %s", songPattern), true
	}

	var insert func(string, int)
	if pos, ok := parseIntParam(args, 1); ok {
		insert = func(s string, i int) { ss.playlist.InsertSong(s, pos+i) }
	} else {
		insert = func(s string, _ int) { ss.playlist.AddSong(s) }
	}

	for i, s := range songs {
		insert(s, i)
	}
	return fmt.Sprintf("%d song(s) added to playlist: %s", len(songs), strings.Join(songs, ", ")), true
}

func (ss *serverState) queueCommand() ssh.Command {
	return ssh.Command{
		Name:     "queue",
		Usage:    "filename [position in playlist]",
		Info:     "adds a song to the playlist",
		ExecFunc: ss.queueCommandExec,
		OptionsFunc: func(prefix string, arg int) []string {
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

func (ss *serverState) playlistCommandExc([]string) (string, bool) {
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
}

func (ss *serverState) playlistCommand() ssh.Command {
	return ssh.Command{
		Name:     "playlist",
		Usage:    "",
		Info:     "prints the current playlist",
		ExecFunc: ss.playlistCommandExc,
	}
}

func (ss *serverState) removeCommand() ssh.Command {
	return ssh.Command{
		Name:  "remove",
		Usage: "position",
		Info:  "removes a song from the playlist",
		ExecFunc: func(args []string) (string, bool) {
			pos, ok := parseIntParam(args, 0)
			if !ok {
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
		ExecFunc: func(args []string) (string, bool) {
			pos, ok := parseIntParam(args, 0)
			if !ok {
				return "", false
			}
			ss.playlist.SetPos(pos)
			return fmt.Sprintf("jumped to %d", pos), true
		},
	}
}

func (ss *serverState) volumeCommand() ssh.Command {
	return ssh.Command{
		Name:  "volume",
		Usage: "volume",
		Info:  "set the playback volume",
		ExecFunc: func(args []string) (string, bool) {
			var ok bool
			ss.volume, ok = parseFloatParam(args, 0)
			if !ok {
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
	return ss.playbackSetCommand(false)
}

func (ss *serverState) resumeCommand() ssh.Command {
	return ss.playbackSetCommand(true)
}

func (ss *serverState) playbackSetCommand(targetPlaying bool) ssh.Command {
	var action string
	if targetPlaying {
		action = "resume"
	} else {
		action = "pause"
	}
	return ssh.Command{
		Name:  action,
		Usage: "",
		Info:  action + "s playback",
		ExecFunc: func([]string) (string, bool) {
			ss.playlist.SetPlaying(targetPlaying)
			return "playback " + action + "d", true
		},
	}
}
