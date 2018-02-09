package schedule

import (
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/ssh"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

// TODO: test queue exec/options

const pathSeparator = string(os.PathSeparator)

type fakeSender struct {
	lastMessage proto.Message
}

func (fs *fakeSender) SendMessage(message proto.Message) error {
	fs.lastMessage = message
	return nil
}

// CommandTester tests a command
type CommandTester interface {
	Test(t *testing.T, command ssh.Command)
}

// OptionsTestCase tests the result of calling the options func on a command
type OptionsTestCase struct {
	Prefix string
	Arg    int
	Result []string
}

func (otc OptionsTestCase) Test(t *testing.T, command ssh.Command) {
	r := command.Options(otc.Prefix, otc.Arg)
	assert.Equal(t, otc.Result, r, "command %s returned wrong options for arg %d with prefix %s", command.Name, otc.Arg, otc.Prefix)
}

// ExecTestCase tests the result of calling the exec func on a command
type ExecTestCase struct {
	Args    []string
	Result  string
	Success bool
	Before  func()
}

func (etc ExecTestCase) Test(t *testing.T, command ssh.Command) {
	if etc.Before != nil {
		etc.Before()
	}
	r, s := command.Exec(etc.Args)
	if assert.Equal(t, etc.Success, s, "command %s returned wrong success flag for args %v", command.Name, etc.Args) && etc.Success {
		assert.Equal(t, etc.Result, r, "command %s returned wrong result for args %v", command.Name, etc.Args)
	}
}

type CommandTesters struct {
	command ssh.Command
	testers []CommandTester
	before  func()
	after   func()
}

func (c CommandTesters) Test(t *testing.T) {
	if c.before != nil {
		c.before()
	}

	for _, tester := range c.testers {
		tester.Test(t, c.command)
	}

	if c.after != nil {
		c.after()
	}
}

func TestServerState_queueCommand(t *testing.T) {
	ss := serverState{}
	ss.playlist = playback.NewPlaylist(0, []string{}, 0)

	cmd := ss.queueCommand()
	assert.NotNil(t, cmd, "serverState queueCommand is nil")

	ct := CommandTesters{
		command: cmd,
		before: func() {
			playback.AudioDir = path.Join("_queue_test_files")
		},
		testers: []CommandTester{
			OptionsTestCase{Prefix: "", Arg: 1, Result: []string{}},
			OptionsTestCase{Prefix: "song", Arg: 0, Result: []string{"song1.mp3", "song2.mp3", "song3.mp3"}},
			OptionsTestCase{Prefix: "dir3", Arg: 0, Result: []string{"dir3" + pathSeparator + "song1.mp3", "dir3" + pathSeparator + "song2.mp3", "dir3" + pathSeparator + "song3.mp3"}},
			ExecTestCase{Args: []string{}, Result: "", Success: false},
			ExecTestCase{Args: []string{"non-existent.mp3"}, Result: "no song matches the glob pattern non-existent.mp3", Success: true},
			ExecTestCase{Args: []string{"song1.mp3"}, Result: "1 song(s) added to playlist: song1.mp3", Success: true},
			ExecTestCase{Args: []string{"dir1/*"}, Result: "3 song(s) added to playlist: dir1" + pathSeparator + "song1.mp3, dir1" + pathSeparator + "song2.mp3, dir1" + pathSeparator + "song3.mp3", Success: true},
			ExecTestCase{Args: []string{"song2.mp3", "abc"}, Result: "1 song(s) added to playlist: song2.mp3", Success: true},
			ExecTestCase{Args: []string{"song3.mp3", "1"}, Result: "1 song(s) added to playlist: song3.mp3", Success: true},
		},
	}

	ct.Test(t)
	assert.Equal(t,
		[]string{"song1.mp3", "song3.mp3", "dir1" + pathSeparator + "song1.mp3",
			"dir1" + pathSeparator + "song2.mp3", "dir1" + pathSeparator + "song3.mp3", "song2.mp3"},
		ss.playlist.Songs(), "serverState queueCommand did not add the songs properly")

	// TODO: options cases
}

func TestServerState_playlistCommand(t *testing.T) {
	ss := serverState{}
	ss.playlist = playback.NewPlaylist(0, []string{}, 0)
	ss.playlist.SetPlaying(false)

	cmd := ss.playlistCommand()
	assert.NotNil(t, cmd, "serverState playlistCommand is nil")

	ct := CommandTesters{
		command: cmd,
		testers: []CommandTester{
			ExecTestCase{Args: []string{}, Result: "Current Playlist (Paused): Empty\nCurrent Song: None", Success: true},
			ExecTestCase{Args: []string{}, Result: "Current Playlist (Playing): Empty\nCurrent Song: None", Success: true, Before: func() { ss.playlist.SetPlaying(true) }},
			ExecTestCase{Args: []string{}, Result: "Current Playlist (Playing): \n  [0] song-1\nCurrent Song: None", Success: true, Before: func() { ss.playlist.AddSong("song-1") }},
			ExecTestCase{Args: []string{}, Result: "Current Playlist (Playing): \n  [0] song-1\n  [1] song-2\nCurrent Song: None", Success: true, Before: func() { ss.playlist.AddSong("song-2") }},
			ExecTestCase{Args: []string{}, Result: "Current Playlist (Playing): \n  [0] song-1\n  [1] song-2\n  [2] song-3\nCurrent Song: None", Success: true, Before: func() { ss.playlist.AddSong("song-3") }},
		},
	}

	ct.Test(t)
}

func TestServerState_removeCommand(t *testing.T) {
	ss := serverState{}
	ss.playlist = playback.NewPlaylist(0, []string{"song-0", "song-1", "song-2", "song-3", "song-4"}, 0)

	cmd := ss.removeCommand()
	assert.NotNil(t, cmd, "serverState removeCommand is nil")

	ct := CommandTesters{
		command: cmd,
		testers: []CommandTester{
			ExecTestCase{Args: []string{}, Result: "", Success: false},
			ExecTestCase{Args: []string{"abc"}, Result: "", Success: false},
			ExecTestCase{Args: []string{"4"}, Result: "removed song song-4 at position 4 from playlist", Success: true},
			ExecTestCase{Args: []string{"0"}, Result: "removed song song-0 at position 0 from playlist", Success: true},
		},
	}

	ct.Test(t)

	assert.Equal(t, []string{"song-1", "song-2", "song-3"}, ss.playlist.Songs(), "songState removeCommand did not call playlist.RemoveSong properly")
}

func TestServerState_jumpCommand(t *testing.T) {
	fs := &fakeSender{}

	ss := serverState{}
	ss.playlist = playback.NewPlaylist(0, []string{"song-0", "song-1", "song-2", "song-3", "song-4"}, 0)
	ss.sender = fs

	cmd := ss.jumpCommand()
	assert.NotNil(t, cmd, "serverState jumpCommand is nil")

	ct := CommandTesters{
		command: cmd,
		testers: []CommandTester{
			ExecTestCase{Args: []string{}, Result: "", Success: false},
			ExecTestCase{Args: []string{"abc"}, Result: "", Success: false},
			ExecTestCase{Args: []string{"3"}, Result: "jumped to 3", Success: true},
		},
	}

	ct.Test(t)

	assert.Equal(t, 3, ss.playlist.Pos(), "songState jumpCommand did not call playlist.SetPos properly")
}

func TestServerState_volumeCommand(t *testing.T) {
	ss := serverState{}
	ss.volume = 0

	fs := &fakeSender{}
	ss.sender = fs

	cmd := ss.volumeCommand()
	assert.NotNil(t, cmd, "serverState volumeCommand is nil")

	ct := CommandTesters{
		command: cmd,
		testers: []CommandTester{
			ExecTestCase{Args: []string{}, Result: "", Success: false},
			ExecTestCase{Args: []string{"abc"}, Result: "", Success: false},
			ExecTestCase{Args: []string{".5"}, Result: "setting volume to 0.500", Success: true},
		},
	}
	ct.Test(t)

	assert.Equal(t, 0.5, ss.volume, "serverState volumeCommand did not set volume properly")
	switch fs.lastMessage.(type) {
	case *comm.SetVolumeRequest:
		svr := fs.lastMessage.(*comm.SetVolumeRequest)
		assert.Equal(t, 0.5, svr.Volume, "serverState volumeCommand did not send the correct SetVolumeRequest")
	default:
		assert.Fail(t, "serverState volumeCommand did not set a SetVolumeRequest: %v", fs.lastMessage)
	}
}

func TestServerState_pauseCommand(t *testing.T) {
	ss := serverState{}
	ss.playlist = playback.NewPlaylist(0, []string{}, 0)
	ss.playlist.SetPlaying(true)

	cmd := ss.pauseCommand()
	assert.NotNil(t, cmd, "serverState pauseCommand is nil")

	ct := CommandTesters{
		command: cmd,
		testers: []CommandTester{
			ExecTestCase{Args: []string{}, Result: "playback paused", Success: true},
		},
	}

	ct.Test(t)

	assert.False(t, ss.playlist.Playing(), "serverState pauseCommand did not pause playback")
}

func TestServerState_resumeCommand(t *testing.T) {
	ss := serverState{}
	ss.playlist = playback.NewPlaylist(0, []string{}, 0)
	ss.playlist.SetPlaying(false)

	cmd := ss.resumeCommand()
	assert.NotNil(t, cmd, "serverState resumeCommand is nil")

	ct := CommandTesters{
		command: cmd,
		testers: []CommandTester{
			ExecTestCase{Args: []string{}, Result: "playback resumed", Success: true},
		},
	}

	ct.Test(t)

	assert.True(t, ss.playlist.Playing(), "serverState resumeCommand did not resume playback")
}

func TestSServerState_playbackSetCommand(t *testing.T) {
	ss := serverState{}

	pause := ss.pauseCommand()
	assert.Equal(t, "pause", pause.Name, "serverState pauseCommand has the wrong name")
	assert.Equal(t, "", pause.Usage, "serverState pauseCommand has the wrong usage")
	assert.Equal(t, "pauses playback", pause.Info, "serverState pauseCommand has the wrong info")

	resume := ss.resumeCommand()
	assert.Equal(t, "resume", resume.Name, "serverState resumeCommand has the wrong name")
	assert.Equal(t, "", resume.Usage, "serverState resumeCommand has the wrong usage")
	assert.Equal(t, "resumes playback", resume.Info, "serverState resumeCommand has the wrong info")
}
