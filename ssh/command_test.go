package ssh

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/playback"
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

type commandTester interface {
	test(t *testing.T, command Command)
}

type optionsTestCase struct {
	prefix string
	arg    int
	result []string
}

func (otc optionsTestCase) test(t *testing.T, command Command) {
	r := command.Options(otc.prefix, otc.arg)
	assert.Equal(t, otc.result, r, "command %s returned wrong options for arg %d with prefix %s", command.Name, otc.arg, otc.prefix)
}

type execTestCase struct {
	args    []string
	result  string
	success bool
}

func (etc execTestCase) test(t *testing.T, command Command) {
	r, s := command.Exec(etc.args)
	if assert.Equal(t, etc.success, s, "command %s returned wrong success flag for args %v", command.Name, etc.args) && etc.success {
		assert.Equal(t, etc.result, r, "command %s returned wrong result for args %v", command.Name, etc.args)
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
	testers []commandTester
	before  func()
	after   func()
}{
	{
		command: helpCommand,
		testers: []commandTester{
			optionsTestCase{prefix: "", arg: 1, result: []string{}},
			optionsTestCase{prefix: "", arg: 0, result: []string{"clear", "exit", "help", "ls"}},
			optionsTestCase{prefix: "he", arg: 0, result: []string{"clear", "exit", "help"}},
			optionsTestCase{prefix: "ls", arg: 0, result: []string{"clear", "exit", "ls"}},
			execTestCase{args: []string{}, success: true, result: "help            retrieves help for a command\nls              lists all songs in the music (sub) directory\nclear           Clears the terminal\nexit            Closes the connection"},
			execTestCase{args: []string{"exit"}, success: true, result: "exit: Closes the connection\nexit"},
			execTestCase{args: []string{"clear"}, success: true, result: "clear: Clears the terminal\nclear"},
			execTestCase{args: []string{"non-existent"}, success: true, result: "Command 'non-existent' does not exist."},
			execTestCase{args: []string{helpCommand.Name}, success: true, result: helpCommand.Name + ": " + helpCommand.Info + "\n" + helpCommand.usage()},
			execTestCase{args: []string{lsCommand.Name}, success: true, result: lsCommand.Name + ": " + lsCommand.Info + "\n" + lsCommand.usage()},
		},
	},
	{
		before:  func() { playback.AudioDir = "_ls_test_files" },
		command: lsCommand,
		testers: []commandTester{
			optionsTestCase{prefix: "", arg: 1, result: []string{}},
			optionsTestCase{prefix: "", arg: 0, result: []string{addPathSep("dir1"), addPathSep("dir2"), addPathSep("dir3")}},
			optionsTestCase{prefix: "dir1", arg: 0, result: []string{addPathSep("dir1")}},
			optionsTestCase{prefix: addPathSep("dir1"), arg: 0, result: []string{addPathSep("dir1", "subdir1"),
				addPathSep("dir1", "subdir2"),
				addPathSep("dir1", "subdir3")}},
			optionsTestCase{prefix: "dir2", arg: 0, result: []string{addPathSep("dir2")}},
			optionsTestCase{prefix: addPathSep("dir2"), arg: 0, result: []string{}},
			execTestCase{args: []string{addPathSep("non-existent", "directory")}, success: true, result: ""},
			execTestCase{args: []string{addPathSep("dir1", "subdir1")}, success: true, result: strings.Join(songFilesInDir("dir1", "subdir1"), "\n")},
			execTestCase{args: []string{addPathSep("dir1", "subdir2")}, success: true, result: strings.Join(songFilesInDir("dir1", "subdir2"), "\n")},
			execTestCase{args: []string{addPathSep("dir1", "subdir3")}, success: true, result: strings.Join(songFilesInDir("dir1", "subdir3"), "\n")},
			execTestCase{args: []string{addPathSep("dir1")}, success: true, result: strings.Join(songFilesInDir("dir1"), "\n") + "\n" +
				strings.Join(songFilesInDir("dir1", "subdir1"), "\n") + "\n" +
				strings.Join(songFilesInDir("dir1", "subdir2"), "\n") + "\n" +
				strings.Join(songFilesInDir("dir1", "subdir3"), "\n")},
			execTestCase{args: []string{addPathSep("dir2")}, success: true, result: strings.Join(songFilesInDir("dir2"), "\n")},
			execTestCase{args: []string{addPathSep("dir3")}, success: true, result: strings.Join(songFilesInDir("dir3"), "\n")},
			execTestCase{args: []string{addPathSep("")}, result: allFilesLsResult},
			execTestCase{args: []string{""}, success: true, result: allFilesLsResult},
			execTestCase{args: []string{}, success: true, result: allFilesLsResult},
		},
	},
}

func TestCommand_OptionsAndExec(t *testing.T) {
	for _, testers := range commandTesters {
		if testers.before != nil {
			testers.before()
		}
		for _, tester := range testers.testers {
			tester.test(t, testers.command)
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
