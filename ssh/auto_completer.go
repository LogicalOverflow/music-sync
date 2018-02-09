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
	if lastChar == ' ' {
		cmd.Parameters = append(cmd.Parameters, "")
	}

	if len(cmd.Parameters) == 0 && lastChar != ' ' {
		newLine = sac.commandsWithPrefix(cmd, pos, newLine)
	} else {
		newLine = sac.commandArgOptionsWithPrefix(cmd, pos, newLine)
	}

	length = pos
	return
}

func (sac *sshAutoCompleter) commandNames() []string {
	cs := []string{"clear", "exit"}
	for _, c := range commands {
		cs = append(cs, c.Name)
	}
	return cs
}

func (sac *sshAutoCompleter) commandsWithPrefix(cmd parser.ParsedCommand, pos int, result [][]rune) [][]rune {
	for _, c := range sac.commandNames() {
		if strings.HasPrefix(c, cmd.Command) {
			result = append(result, []rune((c + " ")[pos:]))
		}
	}
	return result
}

func (sac *sshAutoCompleter) commandArgOptionsWithPrefix(cmd parser.ParsedCommand, pos int, result [][]rune) [][]rune {
	c := commandByName(cmd.Command)
	if c != nil {
		argNum := len(cmd.Parameters) - 1
		arg := cmd.Parameters[argNum]
		for _, o := range sac.filterByPrefix(c.Options(arg, argNum), arg) {
			cmd.Parameters[argNum] = o
			result = append(result, []rune(cmd.Unparse())[pos:])
		}
		cmd.Parameters[argNum] = arg
	}
	return result
}

func (sac *sshAutoCompleter) filterByPrefix(l []string, prefix string) []string {
	r := make([]string, 0, len(l))
	for _, e := range l {
		if strings.HasPrefix(e, prefix) {
			r = append(r, e)
		}
	}
	return r
}
