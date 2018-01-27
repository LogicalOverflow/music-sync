package log

import (
	"strings"
)

type Level int

//noinspection GoUnusedConst
const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) FullName() string {
	return LevelNames[l].Full
}

func (l Level) ShortName() string {
	return LevelNames[l].Short
}

var LevelNames = map[Level]LevelName{
	LevelTrace: {"TRACE", "TRC"},
	LevelDebug: {"DEBUG", "DBG"},
	LevelInfo:  {"INFO", "INF"},
	LevelWarn:  {"WARN", "WRN"},
	LevelError: {"ERROR", "ERR"},
	LevelFatal: {"FATAL", "FTL"},
}

type LevelName struct {
	Full  string
	Short string
}

func shouldLog(level, cutoff Level) bool {
	return cutoff <= level
}

func LevelByName(name string) Level {
	name = strings.ToLower(name)
	for l, n := range LevelNames {
		if strings.ToLower(n.Full) == name || strings.ToLower(n.Short) == name {
			return l
		}
	}
	return LevelOff
}
