package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevelByName(t *testing.T) {
	cases := []struct {
		name  string
		level Level
	}{
		{name: "tRc", level: LevelTrace},
		{name: "dbg", level: LevelDebug},
		{name: "INF", level: LevelInfo},
		{name: "wArN", level: LevelWarn},
		{name: "error", level: LevelError},
		{name: "FATAL", level: LevelFatal},
		{name: "", level: LevelOff},
		{name: "off", level: LevelOff},
		{name: "no-real-level", level: LevelOff},
	}
	for _, c := range cases {
		actual := LevelByName(c.name)
		assert.Equal(t, c.level, actual, "LevelByName(%s) returned the wrong level", c.name)
	}
}
