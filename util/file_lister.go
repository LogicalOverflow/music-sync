// Package util contains utility types and functions
package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilterSongs removes all non-songs from the files slice
func FilterSongs(files []string) []string {
	result := make([]string, 0)
	for _, f := range files {
		if strings.HasSuffix(f, ".mp3") {
			result = append(result, f)
		}
	}
	return result
}

// ListAllFiles recursively lists all songs (files) in songsDir/subDir
func ListAllFiles(songsDir string, subDir string) []string {
	songs := make([]string, 0)
	dir := songsDir
	if subDir != "" {
		dir = filepath.Join(songsDir, subDir)
	}
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
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
		if err != nil {
			return err
		}
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

// ListGlobFiles lists all files in a directory matching the provided glob pattern
func ListGlobFiles(dir, pattern string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return []string{}, err
	}
	if matches == nil || len(matches) == 0 {
		return []string{}, nil
	}

	result := make([]string, 0, len(matches))
	for _, m := range matches {
		if !IsFile(m) {
			continue
		}
		if strings.HasPrefix(m, dir) {
			m = m[len(dir):]
		}
		if strings.HasPrefix(m, "/") || strings.HasPrefix(m, "\\") {
			m = m[1:]
		}
		result = append(result, m)
	}
	return result, nil
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

// IsDir returns true if the path points to an existing directory
func IsDir(p string) bool {
	if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
		return false
	}
	return true
}

// IsFile returns true if the path points to an existing file
func IsFile(p string) bool {
	if fi, err := os.Stat(p); err != nil || fi.IsDir() {
		return false
	}
	return true
}
