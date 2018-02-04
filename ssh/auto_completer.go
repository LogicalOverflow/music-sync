package ssh

import (
	"github.com/LogicalOverflow/music-sync/ssh/parser"
	"strings"
)

type sshAutoCompleter struct{}

func (sac *sshAutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	newLine = make([][]rune, 0)

	lastChar := '\x00'
	if len(line) != 0 {
		lastChar = line[len(line)-1]
	}

	cmd := parser.ParseCommand(string(line))

	if len(cmd.Parameters) == 0 && lastChar != ' ' {
		for _, c := range commands {
			if strings.HasPrefix(c.Name, cmd.Command) {
				newLine = append(newLine, []rune(c.Name + " ")[pos:])
			}
		}

		for _, c := range []string{"clear", "exit"} {
			if strings.HasPrefix(c, cmd.Command) {
				newLine = append(newLine, []rune(c + " ")[pos:])
			}
		}
	} else {
		for _, c := range commands {
			if c.Name == cmd.Command {
				if c.Options != nil {
					if len(cmd.Parameters) == 0 {
						cmd.Parameters = []string{""}
					}
					argNum := len(cmd.Parameters) - 1
					arg := cmd.Parameters[argNum]
					for _, o := range c.Options(arg, argNum) {
						cmd.Parameters[argNum] = o
						newLine = append(newLine, []rune(cmd.Unparse())[pos:])
					}
					cmd.Parameters[argNum] = arg
				}
			}
		}
	}
	length = pos
	return
}
