package parser

import (
	"github.com/LogicalOverflow/music-sync/ssh/parser/shell_parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"reflect"
	"strings"
)

func getParser(line string) *shell_parser.ShellParser {
	inputStream := antlr.NewInputStream(line)
	lexer := shell_parser.NewShellLexer(inputStream)
	lexer.SetChannel(0)
	parser := shell_parser.NewShellParser(antlr.NewCommonTokenStream(lexer, 0))
	return parser
}

// ParsedCommand represents a command that has been parsed
type ParsedCommand struct {
	Command    string
	Parameters []string
}

// Unparse converts a parsed command back to a shell command
func (c ParsedCommand) Unparse() string {
	if len(c.Parameters) == 0 {
		return strings.Replace(c.Command, " ", "\\ ", -1)
	}
	unparsedParams := make([]string, len(c.Parameters))
	for i, p := range c.Parameters {
		unparsedParams[i] = strings.Replace(p, " ", "\\ ", -1)
	}

	return strings.Replace(c.Command, " ", "\\ ", -1) + " " + strings.Join(unparsedParams, " ")
}

func ParseCommand(line string) ParsedCommand {
	p := getParser(line)
	tree := p.Line()
	listener := newShellListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)
	return *listener.cmd
}

type shellListener struct {
	*shell_parser.BaseShellListener
	commandStrings map[int]string
	currentString  string
	cmd            *ParsedCommand
}

func newShellListener() *shellListener {
	l := new(shellListener)
	l.commandStrings = make(map[int]string)
	l.cmd = &ParsedCommand{Command: "", Parameters: make([]string, 0)}
	return l
}

func (s *shellListener) EnterCommandString(ctx *shell_parser.CommandStringContext) {
	s.currentString = ""
}

func (s *shellListener) ExitCharacter(ctx *shell_parser.CharacterContext) {
	if len(ctx.GetText()) != 0 {
		s.currentString += string(ctx.GetText()[len(ctx.GetText())-1])
	}
}

func (s *shellListener) ExitCommandString(ctx *shell_parser.CommandStringContext) {
	s.commandStrings[ctx.GetSourceInterval().Start] = s.currentString
}

func (s *shellListener) ExitCommand(ctx *shell_parser.CommandContext) {
	cctx := ctx.GetChildOfType(0, reflect.TypeOf(&shell_parser.CommandStringContext{}))
	if cctx == nil {
		return
	}
	s.cmd.Command = s.commandStrings[cctx.GetSourceInterval().Start]
}

func (s *shellListener) ExitParameter(ctx *shell_parser.ParameterContext) {
	cctx := ctx.GetChildOfType(0, reflect.TypeOf(&shell_parser.CommandStringContext{}))
	if cctx == nil {
		return
	}
	str := s.commandStrings[cctx.GetSourceInterval().Start]
	s.cmd.Parameters = append(s.cmd.Parameters, str)
}
