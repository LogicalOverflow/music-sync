package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ListAllSongs(songsDir string, subDir string) ([]string) {
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
