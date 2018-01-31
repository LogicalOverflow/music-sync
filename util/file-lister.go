// Package util contains utility types and functions
package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ListAllSongs recursively lists all songs (files) in songsDir/subDir
func ListAllSongs(songsDir string, subDir string) []string {
	songs := make([]string, 0)
	dir := songsDir
	if subDir != "" {
		dir = filepath.Join(songsDir, subDir)
	}
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.HasPrefix(path, songsDir) {
			path = path[len(songsDir):]
		}
		if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\") {
			path = path[1:]
		}
		if path != "" && !f.IsDir() {
			songs = append(songs, path)
		}
		return nil
	})
	return songs
}

// ListAllSubDirs recursively lists all directories in dir
func ListAllSubDirs(dir string) []string {
	dirs := make([]string, 0)
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.HasPrefix(path, dir) {
			path = path[len(dir):]
		}
		if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\") {
			path = path[1:]
		}
		if path != "" && f.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs
}

func ListGlobSongs(dir, pattern string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return []string{}, err
	}
	if matches == nil || len(matches) == 0 {
		return []string{}, nil
	}

	for i := range matches {
		if strings.HasPrefix(matches[i], dir) {
			matches[i] = matches[i][len(dir):]
		}
		if strings.HasPrefix(matches[i], "/") || strings.HasPrefix(matches[i], "\\") {
			matches[i] = matches[i][1:]
		}
	}
	return matches, nil
}

// CheckDir checks whether a directory exists and is in fact a directory (returns an error if that is not the case)
func CheckDir(p string) error {
	fi, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("failed to access %s: %v", p, err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", p)
	}
	return nil
}

// CheckFile checks whether a file exists and is in fact a file (returns an error if that is not the case)
func CheckFile(p string) error {
	fi, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("failed to access %s: %v", p, err)
	}
	if fi.IsDir() {
		return fmt.Errorf("%s is directory", p)
	}
	return nil
}
