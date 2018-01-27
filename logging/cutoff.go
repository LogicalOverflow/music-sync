package log

var CutoffLevels = make(map[string]Level)
var DefaultCutoffLevel = LevelInfo

func cutoffLevel(name string) Level {
	if l, ok := CutoffLevels[name]; ok {
		return l
	} else {
		return DefaultCutoffLevel
	}
}
