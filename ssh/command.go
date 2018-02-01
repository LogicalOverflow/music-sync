package ssh

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/util"
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
	return "Usage: " + command.Name + " " + command.Usage
}

func init() {
	RegisterCommand(Command{
		Name:  "help",
		Usage: "[command name]",
		Info:  "retrieves help for a command",
		Exec: func(args []string) (string, bool) {
			// TODO: add clear/exit
			if len(args) == 0 {
				usages := make([]string, 0, len(commands))
				for _, c := range commands {
					usages = append(usages, fmt.Sprintf("%-15s %s", c.Name, c.Info))
				}
				return strings.Join(usages, "\n"), true
			}

			var target *Command
			for _, c := range commands {
				if c.Name == args[0] {
					target = &c
				}
			}
			if target == nil {
				return fmt.Sprintf("Command %s does not exist.", args[0]), true
			}
			return target.Name + ": " + target.Info + "\n" + target.usage(), true
		},
		Options: func(prefix string, arg int) []string {
			if arg != 0 {
				return []string{}
			}
			options := make([]string, 0)
			for _, c := range commands {
				if strings.HasSuffix(c.Name, prefix) {
					options = append(options, c.Name)
				}
			}
			return options
		},
	})
	RegisterCommand(Command{
		Name:  "ls",
		Usage: "[sub directory]",
		Info:  "lists all files in the music (sub) directory",
		Exec: func(args []string) (string, bool) {
			subDir := ""
			if 0 < len(args) {
				subDir = args[0]
			}
			songs := util.ListAllSongs(playback.AudioDir, subDir)
			return strings.Join(songs, "\n"), true
		},
		Options: func(prefix string, arg int) []string {
			if arg != 0 {
				return []string{}
			}
			songs := util.ListAllSubDirs(playback.AudioDir)
			options := make([]string, 0, len(songs))
			for _, song := range songs {
				if strings.HasPrefix(song, prefix) {
					options = append(options, song)
				}
			}
			return options
		},
	})
}
