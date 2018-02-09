package testutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type commandInterface interface {
	GetName() string
	Exec([]string) (string, bool)
	Options(string, int) []string
}

// CommandTester tests a command
type CommandTester interface {
	Test(t *testing.T, command commandInterface)
}

// OptionsTestCase tests the result of calling the options func on a command
type OptionsTestCase struct {
	Prefix string
	Arg    int
	Result []string
}

// Test executes the OptionsTestCase and asserts the results
func (otc OptionsTestCase) Test(t *testing.T, command commandInterface) {
	r := command.Options(otc.Prefix, otc.Arg)
	assert.Equal(t, otc.Result, r, "command %s returned wrong options for arg %d with prefix %s", command.GetName(), otc.Arg, otc.Prefix)
}

// ExecTestCase tests the result of calling the exec func on a command
type ExecTestCase struct {
	Args    []string
	Result  string
	Success bool
	Before  func()
}

// Test executes the ExecTestCase and asserts the results
func (etc ExecTestCase) Test(t *testing.T, command commandInterface) {
	if etc.Before != nil {
		etc.Before()
	}
	r, s := command.Exec(etc.Args)
	if assert.Equal(t, etc.Success, s, "command %s returned wrong success flag for args %v", command.GetName(), etc.Args) && etc.Success {
		assert.Equal(t, etc.Result, r, "command %s returned wrong result for args %v", command.GetName(), etc.Args)
	}
}

// CommandTesters holds a command and its testers
type CommandTesters struct {
	Command commandInterface
	Testers []CommandTester
	Before  func()
	After   func()
}

// Test executes all testers
func (c CommandTesters) Test(t *testing.T) {
	if c.Before != nil {
		c.Before()
	}

	for _, tester := range c.Testers {
		tester.Test(t, c.Command)
	}

	if c.After != nil {
		c.After()
	}
}
