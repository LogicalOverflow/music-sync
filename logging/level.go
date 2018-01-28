package log

import (
	"strings"
)

// Level is a log level
type Level int

//noinspection GoUnusedConst
// The different log levels
const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOff
)

// FullName returns the full name of a log level
func (l Level) FullName() string {
	return LevelNames[l].Full
}

// ShortName returns the short name of a log level
func (l Level) ShortName() string {
	return LevelNames[l].Short
}

// LevelNames contains the full and short names of the log levels
var LevelNames = map[Level]LevelName{
	LevelTrace: {"TRACE", "TRC"},
	LevelDebug: {"DEBUG", "DBG"},
	LevelInfo:  {"INFO", "INF"},
	LevelWarn:  {"WARN", "WRN"},
	LevelError: {"ERROR", "ERR"},
	LevelFatal: {"FATAL", "FTL"},
}

// LevelName holds the full and short name of a log level
type LevelName struct {
	Full  string
	Short string
}

func shouldLog(level, cutoff Level) bool {
	return cutoff <= level
}

// LevelByName retrieves the level with the given name (defaults to LevelOff)
func LevelByName(name string) Level {
	name = strings.ToLower(name)
	for l, n := range LevelNames {
		if strings.ToLower(n.Full) == name || strings.ToLower(n.Short) == name {
			return l
		}
	}
	return LevelOff
}
