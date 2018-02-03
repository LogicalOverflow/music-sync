package ssh

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRegisterCommand(t *testing.T) {
	oldCommands := make([]Command, len(commands))
	copy(oldCommands, commands)
	initialCommandsLength := len(commands)
	for i := 0; i < 16; i++ {
		c := Command{Name: fmt.Sprintf("test-command-%02d", i)}
		RegisterCommand(c)
		assert.Equal(t, i+initialCommandsLength+1, len(commands), "after registering %d commands, command length is incorrect", i+1)
		for j := 0; j < len(commands)-1; j++ {
			assert.True(t, strings.Compare(commands[j].Name, commands[j+1].Name) <= 0, "after registering %d commands, commands are not ordered by name at index %d", i+1, j)
		}
	}
	commands = make([]Command, len(oldCommands))
	copy(commands, oldCommands)
}
