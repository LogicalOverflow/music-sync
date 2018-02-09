package ssh

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

// CommandTester tests a command
type CommandTester interface {
	Test(t *testing.T, command Command)
}

// OptionsTestCase tests the result of calling the options func on a command
type OptionsTestCase struct {
	Prefix string
	Arg    int
	Result []string
}

func (otc OptionsTestCase) Test(t *testing.T, command Command) {
	r := command.Options(otc.Prefix, otc.Arg)
	assert.Equal(t, otc.Result, r, "command %s returned wrong options for arg %d with prefix %s", command.Name, otc.Arg, otc.Prefix)
}

// ExecTestCase tests the result of calling the exec func on a command
type ExecTestCase struct {
	Args    []string
	Result  string
	Success bool
}

func (etc ExecTestCase) Test(t *testing.T, command Command) {
	r, s := command.Exec(etc.Args)
	if assert.Equal(t, etc.Success, s, "command %s returned wrong success flag for args %v", command.Name, etc.Args) && etc.Success {
		assert.Equal(t, etc.Result, r, "command %s returned wrong result for args %v", command.Name, etc.Args)
	}
}

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

var commandTesters = []struct {
	command Command
	testers []CommandTester
	before  func()
	after   func()
}{
	{
		command: helpCommand,
		testers: []CommandTester{
			OptionsTestCase{Prefix: "", Arg: 1, Result: []string{}},
			OptionsTestCase{Prefix: "", Arg: 0, Result: []string{"clear", "exit", "help", "ls"}},
			OptionsTestCase{Prefix: "he", Arg: 0, Result: []string{"clear", "exit", "help"}},
			OptionsTestCase{Prefix: "ls", Arg: 0, Result: []string{"clear", "exit", "ls"}},
			ExecTestCase{Args: []string{}, Success: true, Result: "help            retrieves help for a command\nls              lists all songs in the music (sub) directory\nclear           Clears the terminal\nexit            Closes the connection"},
			ExecTestCase{Args: []string{"exit"}, Success: true, Result: "exit: Closes the connection\nexit"},
			ExecTestCase{Args: []string{"clear"}, Success: true, Result: "clear: Clears the terminal\nclear"},
			ExecTestCase{Args: []string{"non-existent"}, Success: true, Result: "Command 'non-existent' does not exist."},
			ExecTestCase{Args: []string{helpCommand.Name}, Success: true, Result: helpCommand.Name + ": " + helpCommand.Info + "\n" + helpCommand.usage()},
			ExecTestCase{Args: []string{lsCommand.Name}, Success: true, Result: lsCommand.Name + ": " + lsCommand.Info + "\n" + lsCommand.usage()},
		},
	},
	{
		before:  func() { playback.AudioDir = "_ls_test_files" },
		command: lsCommand,
		testers: []CommandTester{
			OptionsTestCase{Prefix: "", Arg: 1, Result: []string{}},
			OptionsTestCase{Prefix: "", Arg: 0, Result: []string{addPathSep("dir1"), addPathSep("dir2"), addPathSep("dir3")}},
			OptionsTestCase{Prefix: "dir1", Arg: 0, Result: []string{addPathSep("dir1")}},
			OptionsTestCase{Prefix: addPathSep("dir1"), Arg: 0, Result: []string{addPathSep("dir1", "subdir1"),
				addPathSep("dir1", "subdir2"),
				addPathSep("dir1", "subdir3")}},
			OptionsTestCase{Prefix: "dir2", Arg: 0, Result: []string{addPathSep("dir2")}},
			OptionsTestCase{Prefix: addPathSep("dir2"), Arg: 0, Result: []string{}},
			ExecTestCase{Args: []string{addPathSep("non-existent", "directory")}, Success: true, Result: ""},
			ExecTestCase{Args: []string{addPathSep("dir1", "subdir1")}, Success: true, Result: strings.Join(songFilesInDir("dir1", "subdir1"), "\n")},
			ExecTestCase{Args: []string{addPathSep("dir1", "subdir2")}, Success: true, Result: strings.Join(songFilesInDir("dir1", "subdir2"), "\n")},
			ExecTestCase{Args: []string{addPathSep("dir1", "subdir3")}, Success: true, Result: strings.Join(songFilesInDir("dir1", "subdir3"), "\n")},
			ExecTestCase{Args: []string{addPathSep("dir1")}, Success: true, Result: strings.Join(songFilesInDir("dir1"), "\n") + "\n" +
				strings.Join(songFilesInDir("dir1", "subdir1"), "\n") + "\n" +
				strings.Join(songFilesInDir("dir1", "subdir2"), "\n") + "\n" +
				strings.Join(songFilesInDir("dir1", "subdir3"), "\n")},
			ExecTestCase{Args: []string{addPathSep("dir2")}, Success: true, Result: strings.Join(songFilesInDir("dir2"), "\n")},
			ExecTestCase{Args: []string{addPathSep("dir3")}, Success: true, Result: strings.Join(songFilesInDir("dir3"), "\n")},
			ExecTestCase{Args: []string{addPathSep("")}, Success: true, Result: allFilesLsResult},
			ExecTestCase{Args: []string{""}, Success: true, Result: allFilesLsResult},
			ExecTestCase{Args: []string{}, Success: true, Result: allFilesLsResult},
		},
	},
}

func TestCommand_OptionsAndExec(t *testing.T) {
	for _, testers := range commandTesters {
		if testers.before != nil {
			testers.before()
		}
		for _, tester := range testers.testers {
			tester.Test(t, testers.command)
		}
		if testers.after != nil {
			testers.after()
		}
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
