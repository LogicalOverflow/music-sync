package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	var lastMsg string
	var lastLevel Level
	testLogger := Logger(func(msg string, level Level) { lastMsg, lastLevel = msg, level })

	testLogger.Tracef("tracef %s", "tracef")
	assert.Equal(t, "tracef tracef", lastMsg, "tracef resulted in wrong message")
	assert.Equal(t, LevelTrace, lastLevel, "tracef resulted in wrong level")

	testLogger.Debugf("debugf %s", "debugf")
	assert.Equal(t, "debugf debugf", lastMsg, "Debugf resulted in wrong message")
	assert.Equal(t, LevelDebug, lastLevel, "Debugf resulted in wrong level")

	testLogger.Infof("infof %s", "infof")
	assert.Equal(t, "infof infof", lastMsg, "Infof resulted in wrong message")
	assert.Equal(t, LevelInfo, lastLevel, "Infof resulted in wrong level")

	testLogger.Warnf("warnf %s", "warnf")
	assert.Equal(t, "warnf warnf", lastMsg, "Warnf resulted in wrong message")
	assert.Equal(t, LevelWarn, lastLevel, "Warnf resulted in wrong level")

	testLogger.Errorf("errorf %s", "errorf")
	assert.Equal(t, "errorf errorf", lastMsg, "Errorf resulted in wrong message")
	assert.Equal(t, LevelError, lastLevel, "Errorf resulted in wrong level")

	testLogger.Fatalf("fatalf %s", "fatalf")
	assert.Equal(t, "fatalf fatalf", lastMsg, "Fatalf resulted in wrong message")
	assert.Equal(t, LevelFatal, lastLevel, "Fatalf resulted in wrong level")
}
