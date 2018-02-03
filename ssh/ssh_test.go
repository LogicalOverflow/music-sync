package ssh

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

type autoCompleterTestCase struct {
	line    string
	newLine []string
}

func TestSshAutoCompleter_Do(t *testing.T) {
	cases := []autoCompleterTestCase{
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
	}

	autoCompleter := sshAutoCompleter{}

	for _, c := range cases {
		nl, _ := autoCompleter.Do([]rune(c.line), len(c.line))
		resultNL := make([]string, len(nl))
		for i, l := range nl {
			resultNL[i] = string(l)
		}
		if assert.Equal(t, len(c.newLine), len(resultNL), "autoCompleter returned wrong number of new lines (%v vs %v)", c.newLine, resultNL) {
			sort.Strings(c.newLine)
			sort.Strings(resultNL)
			for i := range c.newLine {
				assert.Equal(t, c.newLine[i], resultNL[i], "autoCompleter returned wrong new line (%v vs %v)", c.newLine, nl)
			}
		}
	}
}
