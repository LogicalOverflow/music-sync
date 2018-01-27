package log

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const Format = "[{lvl}] [{name}] {date}: {message}"
const DateFormat = "2006-01-02 15:04:05.000"

func GetLogger(name string) Logger {
	printName := (name + "    ")[:4]
	return func(msg string, level Level) {
		if !shouldLog(level, cutoffLevel(name)) {
			return
		}
		r := strings.NewReplacer("{lvl}", level.ShortName(), "{level}", level.FullName(),
			"{name}", printName, "{date}", time.Now().Format(DateFormat),
			"{message}", msg, "{msg}", msg)
		m := r.Replace(Format)
		if level < LevelWarn {
			fmt.Fprintln(os.Stdout, m)
		} else {
			fmt.Fprintln(os.Stderr, m)
		}
	}
}

type Logger func(msg string, level Level)

func (l Logger) Tracef(format string, a ...interface{}) {
	l.Trace(fmt.Sprintf(format, a...))
}

func (l Logger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l Logger) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l Logger) Warnf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l Logger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l Logger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}

func (l Logger) Trace(msg string) {
	l(msg, LevelTrace)
}

func (l Logger) Debug(msg string) {
	l(msg, LevelDebug)
}

func (l Logger) Info(msg string) {
	l(msg, LevelInfo)
}

func (l Logger) Warn(msg string) {
	l(msg, LevelWarn)
}

func (l Logger) Error(msg string) {
	l(msg, LevelError)
}

func (l Logger) Fatal(msg string) {
	l(msg, LevelFatal)
}
