package cmd

var (
	Commit  = "" // set by build flag
	Tag     = "" // set by build flag
	Version = version()
	Author  = "Leon Vack"
	Email   = ""
)

func version() string {
	v := Tag
	if Tag == "" {
		v = "untagged"
	}
	if Commit != "" {
		c := Commit
		if 8 < len(c) {
			c = c[:8]
		}
		v += "-" + c
	}
	return v
}
