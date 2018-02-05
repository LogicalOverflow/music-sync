package parser

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

var parseCases = []struct {
	line   string
	result ParsedCommand
}{
	{
		line:   "abc def ghi",
		result: ParsedCommand{Command: "abc", Parameters: []string{"def", "ghi"}},
	},
	{
		line:   "abc \"def ghi\"",
		result: ParsedCommand{Command: "abc", Parameters: []string{"\"def", "ghi\""}},
	},
	{
		line:   "abc 'def ghi'",
		result: ParsedCommand{Command: "abc", Parameters: []string{"'def", "ghi'"}},
	},
	{
		line:   "abc def\\ ghi",
		result: ParsedCommand{Command: "abc", Parameters: []string{"def ghi"}},
	},
	{
		line:   "\"abc def\"",
		result: ParsedCommand{Command: "\"abc", Parameters: []string{"def\""}},
	},
	{
		line:   "'abc def'",
		result: ParsedCommand{Command: "'abc", Parameters: []string{"def'"}},
	},
	{
		line:   "abc\\ def",
		result: ParsedCommand{Command: "abc def", Parameters: []string{}},
	},
	{
		line:   "\\\"",
		result: ParsedCommand{Command: "\\\"", Parameters: []string{}},
	},
	{
		line:   "\\'",
		result: ParsedCommand{Command: "\\'", Parameters: []string{}},
	},
	{
		line:   "'\"'",
		result: ParsedCommand{Command: "'\"'", Parameters: []string{}},
	},
	{
		line:   "\"'\"",
		result: ParsedCommand{Command: "\"'\"", Parameters: []string{}},
	},
	{
		line:   "'\\''",
		result: ParsedCommand{Command: "'\\''", Parameters: []string{}},
	},
	{
		line:   "\"\\\"\"",
		result: ParsedCommand{Command: "\"\\\"\"", Parameters: []string{}},
	},
	{
		line:   "\"\\\"\"",
		result: ParsedCommand{Command: "\"\\\"\"", Parameters: []string{}},
	},
	{
		line:   "abc def ghi ",
		result: ParsedCommand{Command: "abc", Parameters: []string{"def", "ghi"}},
	},
}

var unparseCases = []struct {
	cmd    ParsedCommand
	result string
}{
	{
		cmd:    ParsedCommand{Command: "abc", Parameters: []string{"def", "ghi"}},
		result: "abc def ghi ",
	},
	{
		cmd:    ParsedCommand{Command: "abc", Parameters: []string{"def ghi"}},
		result: "abc def\\ ghi ",
	},
	{
		cmd:    ParsedCommand{Command: "abc def", Parameters: []string{"ghi"}},
		result: "abc\\ def ghi ",
	},
	{
		cmd:    ParsedCommand{Command: "abc def", Parameters: []string{}},
		result: "abc\\ def ",
	},
}

func TestParseCommand(t *testing.T) {
	for _, c := range parseCases {
		parsed := ParseCommand(c.line)
		assert.Equal(t, c.result.Command, parsed.Command, "parsing %s resulted in the wrong command name", c.line)
		sort.Strings(c.result.Parameters)
		sort.Strings(parsed.Parameters)
		assert.Equal(t, c.result.Parameters, parsed.Parameters, "parsing %s resulted in the wrong parameters for case %v", c)
	}
}

func TestParsedCommand_Unparse(t *testing.T) {
	for _, c := range unparseCases {
		unparsed := c.cmd.Unparse()
		assert.Equal(t, c.result, unparsed, "unparsing %v yielded the wrong result", c.cmd)
	}
}
