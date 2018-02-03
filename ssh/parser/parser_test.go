package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCommand(t *testing.T) {
	cases := []struct {
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
	}
	for _, c := range cases {
		parsed := ParseCommand(c.line)
		assert.Equal(t, c.result.Command, parsed.Command, "parsing %s resulted in the wrong command name", c.line)
		if assert.Equal(t, len(c.result.Parameters), len(parsed.Parameters), "parsing %s resulted in the wrong number of parameters", c.line) {
			for i := range c.result.Parameters {
				assert.Equal(t, c.result.Parameters[i], parsed.Parameters[i], "parsing %s resulted in the wrong parameter at index %d", c.line, i)
			}
		}
	}
}

func TestParsedCommand_Unparse(t *testing.T) {
	cases := []struct {
		cmd    ParsedCommand
		result string
	}{
		{
			cmd:    ParsedCommand{Command: "abc", Parameters: []string{"def", "ghi"}},
			result: "abc def ghi",
		},
		{
			cmd:    ParsedCommand{Command: "abc", Parameters: []string{"def ghi"}},
			result: "abc def\\ ghi",
		},
	}

	for _, c := range cases {
		unparsed := c.cmd.Unparse()
		assert.Equal(t, c.result, unparsed, "unparsing %v yielded the wrong result", c.cmd)
	}
}
