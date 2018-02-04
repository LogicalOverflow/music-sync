package log

import (
	"bytes"
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

type loggerOutputTestCase struct {
	message       string
	level         Level
	defaultCutoff Level
	namedCutoff   Level
	expectedOut   string
	expectedErr   string
}

var loggerOutputTestCases = []loggerOutputTestCase{
	{
		message:       "filtered by default cutoff",
		level:         LevelTrace,
		defaultCutoff: LevelInfo,
		namedCutoff:   -1,
		expectedOut:   "",
		expectedErr:   "",
	},
	{
		message:       "filtered by named cutoff",
		level:         LevelInfo,
		defaultCutoff: LevelDebug,
		namedCutoff:   LevelWarn,
		expectedOut:   "",
		expectedErr:   "",
	},
	{
		message:       "allowed by named cutoff",
		level:         LevelTrace,
		defaultCutoff: LevelInfo,
		namedCutoff:   LevelTrace,
		expectedOut:   "[TRC] [test] fake-date: allowed by named cutoff\n",
		expectedErr:   "",
	},
	{
		message:       "allowed by default cutoff",
		level:         LevelInfo,
		defaultCutoff: LevelInfo,
		namedCutoff:   -1,
		expectedOut:   "[INF] [test] fake-date: allowed by default cutoff\n",
		expectedErr:   "",
	},
	{
		message:       "trace message",
		level:         LevelTrace,
		defaultCutoff: LevelTrace,
		namedCutoff:   -1,
		expectedOut:   "[TRC] [test] fake-date: trace message\n",
		expectedErr:   "",
	},
	{
		message:       "debug message",
		level:         LevelDebug,
		defaultCutoff: LevelTrace,
		namedCutoff:   -1,
		expectedOut:   "[DBG] [test] fake-date: debug message\n",
		expectedErr:   "",
	},
	{
		message:       "info message",
		level:         LevelInfo,
		defaultCutoff: LevelTrace,
		namedCutoff:   -1,
		expectedOut:   "[INF] [test] fake-date: info message\n",
		expectedErr:   "",
	},
	{
		message:       "warn message",
		level:         LevelWarn,
		defaultCutoff: LevelTrace,
		namedCutoff:   -1,
		expectedOut:   "",
		expectedErr:   "[WRN] [test] fake-date: warn message\n",
	},
	{
		message:       "error message",
		level:         LevelError,
		defaultCutoff: LevelTrace,
		namedCutoff:   -1,
		expectedOut:   "",
		expectedErr:   "[ERR] [test] fake-date: error message\n",
	},
	{
		message:       "fatal message",
		level:         LevelFatal,
		defaultCutoff: LevelTrace,
		namedCutoff:   -1,
		expectedOut:   "",
		expectedErr:   "[FTL] [test] fake-date: fatal message\n",
	},
	{
		message:       "fatal message filtered by off",
		level:         LevelFatal,
		defaultCutoff: LevelOff,
		namedCutoff:   -1,
		expectedOut:   "",
		expectedErr:   "",
	},
}

func TestLoggerOutput(t *testing.T) {
	loggerName := "test"
	logger := GetLogger(loggerName)
	DateProvider = func() string { return "fake-date" }

	for i, c := range loggerOutputTestCases {
		out := bytes.NewBuffer([]byte{})
		OutputWriter = out
		err := bytes.NewBuffer([]byte{})
		ErrorWriter = err

		DefaultCutoffLevel = c.defaultCutoff
		if _, ok := CutoffLevels[loggerName]; ok {
			delete(CutoffLevels, loggerName)
		}
		if 0 <= c.namedCutoff {
			CutoffLevels[loggerName] = c.namedCutoff
		}

		logger(c.message, c.level)
		assert.Equal(t, c.expectedOut, string(out.Bytes()), "Logging %s with level %s in case %d created wrong output on OutputWriter", c.message, c.level, i)
		assert.Equal(t, c.expectedErr, string(err.Bytes()), "Logging %s with level %s in case %d created wrong output on ErrorWriter", c.message, c.level, i)
	}
}
