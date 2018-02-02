// Package log provides utilities for logging
package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// OutputWriter is the writer to which normal log messages are written (level < LevelWarn)
var OutputWriter io.Writer = os.Stdout

// ErrorWriter is the writer to which Warn/Error/Fatal log messages are written
var ErrorWriter io.Writer = os.Stderr

// Format describes how a log message is formatted
const Format = "[{lvl}] [{name}] {date}: {message}"

// DateFormat describes how the date in a log message is formatted
const DateFormat = "2006-01-02 15:04:05.000"

// DateProvider returns the formatted date as string, use by the logger
var DateProvider = func() string { return time.Now().Format(DateFormat) }

// GetLogger returns a Logger to log from a given package
func GetLogger(name string) Logger {
	printName := (name + "    ")[:4]
	return func(msg string, level Level) {
		if !shouldLog(level, cutoffLevel(name)) {
			return
		}
		r := strings.NewReplacer("{lvl}", level.ShortName(), "{level}", level.FullName(),
			"{name}", printName, "{date}", DateProvider(),
			"{message}", msg, "{msg}", msg)
		m := r.Replace(Format)
		if level < LevelWarn {
			fmt.Fprintln(OutputWriter, m)
		} else {
			fmt.Fprintln(ErrorWriter, m)
		}
	}
}

// Logger is a function to log a single message
type Logger func(msg string, level Level)

// Tracef formats a message and logs it with LevelTrace
func (l Logger) Tracef(format string, a ...interface{}) {
	l.Trace(fmt.Sprintf(format, a...))
}

// Debugf formats a message and logs it with LevelDebug
func (l Logger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

// Infof formats a message and logs it with LevelInfo
func (l Logger) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

// Warnf formats a message and logs it with LevelWarn
func (l Logger) Warnf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}

// Errorf formats a message and logs it with LevelError
func (l Logger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

// Fatalf formats a message and logs it with LevelFatal
func (l Logger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}

// Trace logs a message with LevelTrace
func (l Logger) Trace(msg string) {
	l(msg, LevelTrace)
}

// Debug logs a message with LevelDebug
func (l Logger) Debug(msg string) {
	l(msg, LevelDebug)
}

// Info logs a message with LevelInfo
func (l Logger) Info(msg string) {
	l(msg, LevelInfo)
}

// Warn logs a message with LevelWarn
func (l Logger) Warn(msg string) {
	l(msg, LevelWarn)
}

// Error logs a message with LevelError
func (l Logger) Error(msg string) {
	l(msg, LevelError)
}

// Fatal logs a message with LevelFatal
func (l Logger) Fatal(msg string) {
	l(msg, LevelFatal)
}
