package log

// CutoffLevels maps the loggers by name to their cutoff level
var CutoffLevels = make(map[string]Level)

// DefaultCutoffLevel is used as cutoff level for loggers with a name for whom no cutoff level is defined in CutoffLevels
var DefaultCutoffLevel = LevelInfo

func cutoffLevel(name string) Level {
	if l, ok := CutoffLevels[name]; ok {
		return l
	}
	return DefaultCutoffLevel
}
