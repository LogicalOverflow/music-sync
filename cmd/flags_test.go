package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"testing"
)

func TestAddLoggingFlags(t *testing.T) {
	names := []string{"comm-logging", "play-logging", "shed-logging", "ssh-logging", "time-logging", "logging"}
	f := AddLoggingFlags([]cli.Flag{})
	require.Equal(t, len(names), len(f), "AddLoggingFlags did not add the right number of flags")
	for i := range names {
		assert.Equal(t, names[i], f[i].GetName(), "AddLoggingFlags has the wrong flag at index %d", i)
	}
}
