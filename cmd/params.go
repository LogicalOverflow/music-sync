// Package cmd provides utilities used to create the different command line interfaces
package cmd

// Tag is the current git tag (set by build flag)
var Tag = ""

// Commit is the current git commit (set by build flag)
var Commit = ""

// Version is the current version (Tag-Commit)
var Version = version()

// Author is the author of the app
var Author = "Leon Vack"

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
