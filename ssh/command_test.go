package ssh

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/testutil"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestRegisterCommand(t *testing.T) {
	oldCommands := make([]Command, len(commands))
	copy(oldCommands, commands)
	initialCommandsLength := len(commands)
	for i := 0; i < 16; i++ {
		c := Command{Name: fmt.Sprintf("test-command-%02d", i)}
		RegisterCommand(c)
		assert.Equal(t, i+initialCommandsLength+1, len(commands), "after registering %d commands, command length is incorrect", i+1)
		for j := 0; j < len(commands)-1; j++ {
			assert.True(t, strings.Compare(commands[j].Name, commands[j+1].Name) <= 0, "after registering %d commands, commands are not ordered by name at index %d", i+1, j)
		}
	}
	commands = make([]Command, len(oldCommands))
	copy(commands, oldCommands)
}

func TestCommand_usage(t *testing.T) {
	cases := []struct {
		name   string
		usage  string
		result string
	}{
		{
			name:   "test-name",
			usage:  "test-usage",
			result: "Usage: test-name test-usage",
		},
		{
			name:   "test-name",
			usage:  "",
			result: "No usage information for command 'test-name'",
		},
	}

	for _, c := range cases {
		cmd := Command{Name: c.name, Usage: c.usage}
		assert.Equal(t, c.result, cmd.usage(), "command usage returned wrong value for case %v", c)
	}
}

var allFilesLsResult = strings.Join(songFilesInDir("dir1"), "\n") + "\n" +
	strings.Join(songFilesInDir("dir1", "subdir1"), "\n") + "\n" +
	strings.Join(songFilesInDir("dir1", "subdir2"), "\n") + "\n" +
	strings.Join(songFilesInDir("dir1", "subdir3"), "\n") + "\n" +
	strings.Join(songFilesInDir("dir2"), "\n") + "\n" +
	strings.Join(songFilesInDir("dir3"), "\n") + "\n" +
	strings.Join(songFilesInDir(), "\n")

var commandTesters = []testutil.CommandTesters{
	{
		Command: helpCommand,
		Testers: []testutil.CommandTester{
			testutil.OptionsTestCase{Prefix: "", Arg: 1, Result: []string{}},
			testutil.OptionsTestCase{Prefix: "", Arg: 0, Result: []string{"clear", "exit", "help", "ls"}},
			testutil.OptionsTestCase{Prefix: "he", Arg: 0, Result: []string{"clear", "exit", "help"}},
			testutil.OptionsTestCase{Prefix: "ls", Arg: 0, Result: []string{"clear", "exit", "ls"}},
			testutil.ExecTestCase{Args: []string{}, Success: true, Result: "help            retrieves help for a command\nls              lists all songs in the music (sub) directory\nclear           Clears the terminal\nexit            Closes the connection"},
			testutil.ExecTestCase{Args: []string{"exit"}, Success: true, Result: "exit: Closes the connection\nexit"},
			testutil.ExecTestCase{Args: []string{"clear"}, Success: true, Result: "clear: Clears the terminal\nclear"},
			testutil.ExecTestCase{Args: []string{"non-existent"}, Success: true, Result: "Command 'non-existent' does not exist."},
			testutil.ExecTestCase{Args: []string{helpCommand.Name}, Success: true, Result: helpCommand.Name + ": " + helpCommand.Info + "\n" + helpCommand.usage()},
			testutil.ExecTestCase{Args: []string{lsCommand.Name}, Success: true, Result: lsCommand.Name + ": " + lsCommand.Info + "\n" + lsCommand.usage()},
		},
	},
	{
		Before:  func() { playback.AudioDir = "_ls_test_files" },
		Command: lsCommand,
		Testers: []testutil.CommandTester{
			testutil.OptionsTestCase{Prefix: "", Arg: 1, Result: []string{}},
			testutil.OptionsTestCase{Prefix: "", Arg: 0, Result: []string{addPathSep("dir1"), addPathSep("dir2"), addPathSep("dir3")}},
			testutil.OptionsTestCase{Prefix: "dir1", Arg: 0, Result: []string{addPathSep("dir1")}},
			testutil.OptionsTestCase{Prefix: addPathSep("dir1"), Arg: 0, Result: []string{addPathSep("dir1", "subdir1"),
				addPathSep("dir1", "subdir2"),
				addPathSep("dir1", "subdir3")}},
			testutil.OptionsTestCase{Prefix: "dir2", Arg: 0, Result: []string{addPathSep("dir2")}},
			testutil.OptionsTestCase{Prefix: addPathSep("dir2"), Arg: 0, Result: []string{}},
			testutil.ExecTestCase{Args: []string{addPathSep("non-existent", "directory")}, Success: true, Result: ""},
			testutil.ExecTestCase{Args: []string{addPathSep("dir1", "subdir1")}, Success: true, Result: strings.Join(songFilesInDir("dir1", "subdir1"), "\n")},
			testutil.ExecTestCase{Args: []string{addPathSep("dir1", "subdir2")}, Success: true, Result: strings.Join(songFilesInDir("dir1", "subdir2"), "\n")},
			testutil.ExecTestCase{Args: []string{addPathSep("dir1", "subdir3")}, Success: true, Result: strings.Join(songFilesInDir("dir1", "subdir3"), "\n")},
			testutil.ExecTestCase{Args: []string{addPathSep("dir1")}, Success: true, Result: strings.Join(songFilesInDir("dir1"), "\n") + "\n" +
				strings.Join(songFilesInDir("dir1", "subdir1"), "\n") + "\n" +
				strings.Join(songFilesInDir("dir1", "subdir2"), "\n") + "\n" +
				strings.Join(songFilesInDir("dir1", "subdir3"), "\n")},
			testutil.ExecTestCase{Args: []string{addPathSep("dir2")}, Success: true, Result: strings.Join(songFilesInDir("dir2"), "\n")},
			testutil.ExecTestCase{Args: []string{addPathSep("dir3")}, Success: true, Result: strings.Join(songFilesInDir("dir3"), "\n")},
			testutil.ExecTestCase{Args: []string{addPathSep("")}, Success: true, Result: allFilesLsResult},
			testutil.ExecTestCase{Args: []string{""}, Success: true, Result: allFilesLsResult},
			testutil.ExecTestCase{Args: []string{}, Success: true, Result: allFilesLsResult},
		},
	},
	{
		Command: Command{},
		Testers: []testutil.CommandTester{
			testutil.OptionsTestCase{Prefix: "", Arg: 0, Result: []string{}},
			testutil.OptionsTestCase{Prefix: "", Arg: 1, Result: []string{}},
		},
	},
}

func TestCommand_OptionsAndExec(t *testing.T) {
	for _, testers := range commandTesters {
		testers.Test(t)
	}
}

func addPathSep(p ...string) string {
	return strings.Join(p, string(os.PathSeparator)) + string(os.PathSeparator)
}

func songFilesInDir(p ...string) []string {
	ss := make([]string, 0)
	bp := addPathSep(p...)
	if len(p) == 0 {
		bp = ""
	}
	for _, s := range []string{"song1.mp3", "song2.mp3", "song3.mp3"} {
		ss = append(ss, bp+s)
	}
	return ss
}
