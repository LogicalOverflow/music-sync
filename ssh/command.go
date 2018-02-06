package ssh

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/util"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var commands = make([]Command, 0)

type commandsByName []Command

func (c commandsByName) Len() int               { return len(c) }
func (c commandsByName) Less(i int, j int) bool { return strings.Compare(c[i].Name, c[j].Name) < 0 }
func (c commandsByName) Swap(i int, j int)      { c[i], c[j] = c[j], c[i] }

// RegisterCommand registers a command to allow its use from ssh control interface
func RegisterCommand(c Command) {
	commands = append(commands, c)
	sort.Sort(commandsByName(commands))
}

// Command describes a ssh command
type Command struct {
	Name  string // Name is the name of the command, which is used to access it from the terminal
	Usage string // Usage contains usage information for the command
	Info  string // Info contains information about what the command does
	// Exec runs the command. It is passed the arguments as a string slice.
	// If it returns false, a usage message is printed, otherwise the returned string is printed
	Exec func([]string) (string, bool)
	// Options (optional) is used for auto completion. It is passed a prefix and the number of the argument and should
	// return all possible completion options
	Options func(prefix string, arg int) []string
}

func (command Command) usage() string {
	if command.Usage == "" {
		return fmt.Sprintf("No usage information for command '%s'", command.Name)
	}
	return "Usage: " + command.Name + " " + command.Usage
}

func commandByName(name string) *Command {
	for _, c := range commands {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

var helpCommand = Command{
	Name:  "help",
	Usage: "[command name]",
	Info:  "retrieves help for a command",
	Exec: func(args []string) (string, bool) {
		if len(args) == 0 {
			usages := make([]string, 0, len(commands))
			for _, c := range commands {
				usages = append(usages, fmt.Sprintf("%-15s %s", c.Name, c.Info))
			}
			usages = append(usages,
				fmt.Sprintf("%-15s %s", "clear", "Clears the terminal"),
				fmt.Sprintf("%-15s %s", "exit", "Closes the connection"),
			)

			return strings.Join(usages, "\n"), true
		}

		var target *Command
		for _, c := range commands {
			if c.Name == args[0] {
				target = &Command{}
				*target = c
			}
		}
		if target == nil {
			switch args[0] {
			case "clear":
				return "clear: Clears the terminal\nclear", true
			case "exit":
				return "exit: Closes the connection\nexit", true
			default:
				return fmt.Sprintf("Command '%s' does not exist.", args[0]), true
			}
		}
		return target.Name + ": " + target.Info + "\n" + target.usage(), true
	},
	Options: func(prefix string, arg int) []string {
		if arg != 0 {
			return []string{}
		}
		options := []string{"clear", "exit"}
		for _, c := range commands {
			if strings.HasPrefix(c.Name, prefix) {
				options = append(options, c.Name)
			}
		}
		return options
	},
}

var lsCommand = Command{
	Name:  "ls",
	Usage: "[sub directory]",
	Info:  "lists all songs in the music (sub) directory",
	Exec: func(args []string) (string, bool) {
		subDir := ""
		if 0 < len(args) {
			subDir = args[0]
		}
		songs := util.FilterSongs(util.ListAllFiles(playback.AudioDir, subDir))
		return strings.Join(songs, "\n"), true
	},
	Options: func(prefix string, arg int) []string {
		if arg != 0 {
			return []string{}
		}
		subDirs := util.ListAllSubDirs(playback.AudioDir)
		options := make([]string, 0, len(subDirs))
		for _, subDir := range subDirs {
			if !strings.HasPrefix(subDir, prefix) {
				continue
			}
			if d, f := filepath.Split(subDir[len(prefix):]); d == "" && f == subDir[len(prefix):] {
				options = append(options, subDir+string(os.PathSeparator))
			}
		}
		return options
	},
}

func init() {
	RegisterCommand(helpCommand)
	RegisterCommand(lsCommand)
}
