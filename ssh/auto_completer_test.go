package ssh

import (
	"github.com/LogicalOverflow/music-sync/ssh/parser"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

type autoCompleterTestCase struct {
	line    string
	newLine []string
}

var autoCompleterTestCases = []autoCompleterTestCase{
	{
		line:    "c",
		newLine: []string{"lear "},
	},
	{
		line:    "e",
		newLine: []string{"xit "},
	},
	{
		line:    "",
		newLine: []string{"clear ", "exit ", "help ", "ls "},
	},
	{
		line:    "l",
		newLine: []string{"s "},
	},
	{
		line:    "exit",
		newLine: []string{" "},
	},
	{
		line:    "exit ",
		newLine: []string{},
	},
}

var filterByPrefixCases = []struct {
	strings []string
	prefix  string
	result  []string
}{
	{
		strings: []string{"wrong-string-1", "prefix-string-1", "prefix-string-2", "wrong-string-2"},
		prefix:  "prefix",
		result:  []string{"prefix-string-1", "prefix-string-2"},
	},
	{
		strings: []string{"string-1", "string-2", "string-3"},
		prefix:  "",
		result:  []string{"string-1", "string-2", "string-3"},
	},
	{
		strings: []string{},
		prefix:  "",
		result:  []string{},
	},
}

var commandArgOptionsWithPrefixCases = []struct {
	parameters    []string
	options       []string
	pos           int
	initialResult [][]rune
	result        []string
}{
	{
		parameters:    []string{""},
		options:       []string{},
		pos:           0,
		initialResult: [][]rune{},
		result:        []string{},
	},
	{
		parameters:    []string{"prefix-"},
		options:       []string{},
		pos:           7,
		initialResult: [][]rune{},
		result:        []string{},
	},
	{
		parameters:    []string{"prefix-"},
		options:       []string{"option-1", "prefix-option-2", "option-3", "prefix-option-4"},
		pos:           7,
		initialResult: [][]rune{},
		result:        []string{"option-2 ", "option-4 "},
	},
	{
		parameters:    []string{"param-1", "param-2", "prefix-"},
		options:       []string{"option-1", "prefix-option-2", "option-3", "prefix-option-4"},
		pos:           23,
		initialResult: [][]rune{},
		result:        []string{"option-2 ", "option-4 "},
	},
	{
		parameters:    []string{"prefix "},
		options:       []string{"option 1", "prefix option 2", "option 3", "prefix option 4"},
		pos:           8,
		initialResult: [][]rune{},
		result:        []string{"option\\ 2 ", "option\\ 4 "},
	},
}

func TestSshAutoCompleter_Do(t *testing.T) {
	autoCompleter := sshAutoCompleter{}
	for _, c := range autoCompleterTestCases {
		nl, _ := autoCompleter.Do([]rune(c.line), len(c.line))
		resultNL := make([]string, len(nl))
		for i, l := range nl {
			resultNL[i] = string(l)
		}
		assertStringSliceEqual(t, c.newLine, resultNL, "autoCompleter returned")
	}
}

func assertStringSliceEqual(t *testing.T, expected, actual []string, name string) {
	sort.Strings(expected)
	sort.Strings(actual)
	assert.Equal(t, expected, actual, "%s wrong strings", name)
}

func TestSshAutoCompleter_filterByPrefix(t *testing.T) {
	autoCompleter := sshAutoCompleter{}
	for _, c := range filterByPrefixCases {
		actual := autoCompleter.filterByPrefix(c.strings, c.prefix)
		assert.Equal(t, c.result, actual, "filterByPrefix result is wrong for case %v", c)
	}
}

func TestSshAutoCompleter_commandArgOptionsWithPrefix(t *testing.T) {
	commandsOld := make([]Command, len(commands))
	copy(commandsOld, commands)

	autoCompleter := sshAutoCompleter{}
	var options []string
	testCmd := Command{Name: "test-command", Options: func(string, int) []string { return options }}
	RegisterCommand(testCmd)

	for _, c := range commandArgOptionsWithPrefixCases {
		options = c.options
		cmd := parser.ParsedCommand{Command: testCmd.Name, Parameters: c.parameters}
		c.pos += len(testCmd.Name) + 1
		actualRune := autoCompleter.commandArgOptionsWithPrefix(cmd, c.pos, c.initialResult)
		assert.Equal(t, runeSliceSliceToStringSlice(actualRune), c.result, "commandArgOptionsWithPrefix result is wrong for case %v", c)
	}

	copy(commands, commandsOld)
}

func runeSliceSliceToStringSlice(rss [][]rune) []string {
	ss := make([]string, len(rss))
	for i, rs := range rss {
		ss[i] = string(rs)
	}
	return ss
}
