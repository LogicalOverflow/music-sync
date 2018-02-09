package schedule

import (
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/ssh"
	"github.com/LogicalOverflow/music-sync/testutil"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

const pathSeparator = string(os.PathSeparator)

type fakeSender struct {
	lastMessage proto.Message
}

func (fs *fakeSender) SendMessage(message proto.Message) error {
	fs.lastMessage = message
	return nil
}

var noArgsError = testutil.ExecTestCase{Args: []string{}, Result: "", Success: false}
var firstArgNoNumberError = testutil.ExecTestCase{Args: []string{"not-a-number"}, Result: "", Success: false}

func TestServerState_queueCommand(t *testing.T) {
	ss := newTestServerState([]string{}, false)

	cmd := ss.queueCommand()
	assert.NotNil(t, cmd, "serverState queueCommand is nil")

	ct := testutil.CommandTesters{
		Command: cmd,
		Before: func() {
			playback.AudioDir = path.Join("_queue_test_files")
		},
		Testers: []testutil.CommandTester{
			testutil.OptionsTestCase{Prefix: "", Arg: 1, Result: []string{}},
			testutil.OptionsTestCase{Prefix: "song", Arg: 0, Result: []string{"song1.mp3", "song2.mp3", "song3.mp3"}},
			testutil.OptionsTestCase{Prefix: "dir3", Arg: 0, Result: []string{"dir3" + pathSeparator + "song1.mp3", "dir3" + pathSeparator + "song2.mp3", "dir3" + pathSeparator + "song3.mp3"}},
			noArgsError,
			testutil.ExecTestCase{Args: []string{"non-existent.mp3"}, Result: "no song matches the glob pattern non-existent.mp3", Success: true},
			testutil.ExecTestCase{Args: []string{"song1.mp3"}, Result: "1 song(s) added to playlist: song1.mp3", Success: true},
			testutil.ExecTestCase{Args: []string{"dir1/*"}, Result: "3 song(s) added to playlist: dir1" + pathSeparator + "song1.mp3, dir1" + pathSeparator + "song2.mp3, dir1" + pathSeparator + "song3.mp3", Success: true},
			testutil.ExecTestCase{Args: []string{"song2.mp3", "abc"}, Result: "1 song(s) added to playlist: song2.mp3", Success: true},
			testutil.ExecTestCase{Args: []string{"song3.mp3", "1"}, Result: "1 song(s) added to playlist: song3.mp3", Success: true},
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
	ss := newTestServerState([]string{}, false)

	cmd := ss.playlistCommand()
	assert.NotNil(t, cmd, "serverState playlistCommand is nil")

	ct := testutil.CommandTesters{
		Command: cmd,
		Testers: []testutil.CommandTester{
			testutil.ExecTestCase{Args: []string{}, Result: "Current Playlist (Paused): Empty\nCurrent Song: None", Success: true},
			testutil.ExecTestCase{Args: []string{}, Result: "Current Playlist (Playing): Empty\nCurrent Song: None", Success: true, Before: func() { ss.playlist.SetPlaying(true) }},
			testutil.ExecTestCase{Args: []string{}, Result: "Current Playlist (Playing): \n  [0] song-1\nCurrent Song: None", Success: true, Before: func() { ss.playlist.AddSong("song-1") }},
			testutil.ExecTestCase{Args: []string{}, Result: "Current Playlist (Playing): \n  [0] song-1\n  [1] song-2\nCurrent Song: None", Success: true, Before: func() { ss.playlist.AddSong("song-2") }},
			testutil.ExecTestCase{Args: []string{}, Result: "Current Playlist (Playing): \n  [0] song-1\n  [1] song-2\n  [2] song-3\nCurrent Song: None", Success: true, Before: func() { ss.playlist.AddSong("song-3") }},
		},
	}

	ct.Test(t)
}

func TestServerState_removeCommand(t *testing.T) {
	ss := newTestServerState([]string{"song-0", "song-1", "song-2", "song-3", "song-4"}, false)

	cmd := ss.removeCommand()
	assert.NotNil(t, cmd, "serverState removeCommand is nil")

	ct := testutil.CommandTesters{
		Command: cmd,
		Testers: []testutil.CommandTester{
			noArgsError, firstArgNoNumberError,
			testutil.ExecTestCase{Args: []string{"4"}, Result: "removed song song-4 at position 4 from playlist", Success: true},
			testutil.ExecTestCase{Args: []string{"0"}, Result: "removed song song-0 at position 0 from playlist", Success: true},
		},
	}

	ct.Test(t)

	assert.Equal(t, []string{"song-1", "song-2", "song-3"}, ss.playlist.Songs(), "songState removeCommand did not call playlist.RemoveSong properly")
}

func TestServerState_jumpCommand(t *testing.T) {
	ss := newTestServerState([]string{"song-0", "song-1", "song-2", "song-3", "song-4"}, false)

	fs := &fakeSender{}
	ss.sender = fs

	cmd := ss.jumpCommand()
	assert.NotNil(t, cmd, "serverState jumpCommand is nil")

	ct := testutil.CommandTesters{
		Command: cmd,
		Testers: []testutil.CommandTester{
			noArgsError, firstArgNoNumberError,
			testutil.ExecTestCase{Args: []string{"3"}, Result: "jumped to 3", Success: true},
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

	ct := testutil.CommandTesters{
		Command: cmd,
		Testers: []testutil.CommandTester{
			noArgsError, firstArgNoNumberError,
			testutil.ExecTestCase{Args: []string{".5"}, Result: "setting volume to 0.500", Success: true},
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
	testServerStatePauseOrResumeCommand(t, false)
}

func TestServerState_resumeCommand(t *testing.T) {
	testServerStatePauseOrResumeCommand(t, true)
}

func testServerStatePauseOrResumeCommand(t *testing.T, playing bool) {
	ss := newTestServerState([]string{}, !playing)

	var cmd ssh.Command
	var key string
	if playing {
		key = "resume"
		cmd = ss.resumeCommand()
	} else {
		key = "pause"
		cmd = ss.pauseCommand()
	}
	assert.NotNil(t, cmd, "serverState %sCommand is nil", key)

	ct := testutil.CommandTesters{
		Command: cmd,
		Testers: []testutil.CommandTester{testutil.ExecTestCase{Args: []string{}, Result: "playback " + key + "d", Success: true}},
	}
	ct.Test(t)

	assert.Equal(t, playing, ss.playlist.Playing(), "serverState %sCommand did not resume playback", key)
}

func TestServerState_playbackSetCommand(t *testing.T) {
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

func newTestServerState(songs []string, playing bool) serverState {
	ss := serverState{}
	ss.playlist = playback.NewPlaylist(0, songs, 0)
	ss.playlist.SetPlaying(playing)
	return ss
}
