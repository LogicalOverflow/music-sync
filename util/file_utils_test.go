package util

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path"
	"path/filepath"
	"sort"
	"testing"
)

const testDir = "_file_lister_test_files"

// _file_lister_test_files:
// - dir1
//   - subdir1
//     - file1
//     - file2
//     - file3
//   - file1
//   - file2
//   - file3
// - dir2
//   - file1
//   - file2
//   - file3
// - dir3
//   - file1
//   - file2
//   - file3
// - file1
// - file2
// - file3

func TestListAllFiles(t *testing.T) {
	songsDir1 := ListAllFiles(testDir, "dir1")
	require.Equal(t, 6, len(songsDir1), "number of files found by ListAllFiles with sub dir is wrong")

	sort.Strings(songsDir1)

	expectedFiles := [][]string{
		{"dir1", "file1"}, {"dir1", "file2"}, {"dir1", "file3"},
		{"dir1", "subdir1", "file1"}, {"dir1", "subdir1", "file2"}, {"dir1", "subdir1", "file3"},
	}
	assertFilesAndFolders(t, expectedFiles, songsDir1, "file", "ListAllFiles with sub dir")

	songs := ListAllFiles(testDir, "")
	require.Equal(t, 15, len(songs), "number of files found by ListAllFiles without sub dir is wrong")

	sort.Strings(songs)

	expectedFiles = [][]string{
		{"dir1", "file1"}, {"dir1", "file2"}, {"dir1", "file3"},
		{"dir1", "subdir1", "file1"}, {"dir1", "subdir1", "file2"}, {"dir1", "subdir1", "file3"},
		{"dir2", "file1"}, {"dir2", "file2"}, {"dir2", "file3"},
		{"dir3", "file1"}, {"dir3", "file2"}, {"dir3", "file3"},
		{"file1"}, {"file2"}, {"file3"},
	}
	assertFilesAndFolders(t, expectedFiles, songs, "file", "ListAllFiles without sub dir")
}

func TestListAllSubDirs(t *testing.T) {
	subDirs := ListAllSubDirs(testDir)
	require.Equal(t, 4, len(subDirs), "number of sub dirs found by ListAllSubDirs is wrong")

	sort.Strings(subDirs)

	assert.Equal(t, "dir1", subDirs[0], "dir found by ListAllSubDirs is wrong")
	assert.Equal(t, filepath.Join("dir1", "subdir1"), subDirs[1], "dir found by ListAllSubDirs is wrong")
	assert.Equal(t, "dir2", subDirs[2], "dir found by ListAllSubDirs is wrong")
	assert.Equal(t, "dir3", subDirs[3], "dir found by ListAllSubDirs is wrong")
}

func TestListGlobFiles(t *testing.T) {
	songsDir1, err := ListGlobFiles(testDir, "dir1/*")
	require.Nil(t, err, "ListGlobFiles for dir1/* caused an error")

	sort.Strings(songsDir1)

	expectedFiles := [][]string{{"dir1", "file1"}, {"dir1", "file2"}, {"dir1", "file3"}}
	assertFilesAndFolders(t, expectedFiles, songsDir1, "file", "ListGlobFiles for dir/*")

	songs1Deep, err := ListGlobFiles(testDir, "*/*")
	require.Nil(t, err, "ListGlobFiles for */* caused an error")

	sort.Strings(songs1Deep)

	expectedFiles = [][]string{
		{"dir1", "file1"}, {"dir1", "file2"}, {"dir1", "file3"},
		{"dir2", "file1"}, {"dir2", "file2"}, {"dir2", "file3"},
		{"dir3", "file1"}, {"dir3", "file2"}, {"dir3", "file3"},
	}
	assertFilesAndFolders(t, expectedFiles, songs1Deep, "file", "ListGlobFiles for */*")
}

func assertFilesAndFolders(t *testing.T, expected [][]string, actual []string, name, function string) {
	if assert.Equal(t, len(expected), len(actual), "number of %ss found by %s is wrong", name, function) {
		for i := range actual {
			expectedPath := filepath.Join(expected[i]...)
			assert.Equal(t, expectedPath, actual[i], "%s found by %s is wrong", name, function)
		}
	}
}

func TestFilterSongs(t *testing.T) {
	cases := [][2][]string{
		{{"no-song.json", "a-song.mp3", "another-song.mp3", "another-non-song.bin"}, {"a-song.mp3", "another-song.mp3"}},
		{{"no-song", "nope"}, {}},
		{{".mp3", "this-song.mp3"}, {".mp3", "this-song.mp3"}},
	}
	for _, c := range cases {
		actual := FilterSongs(c[0])
		expected := c[1]
		if assert.Equal(t, len(expected), len(actual), "FilterSong result length is wrong for %v", c[0]) {
			for i := range expected {
				assert.Equal(t, expected[i], actual[i], "FilterSong result at index %d is wrong for %v", i, c[0])
			}
		}
	}
}

type isDirFileCase struct {
	name   string
	isDir  bool
	exists bool
}

var isDirFileCases = []isDirFileCase{
	{name: "non-existent", isDir: false, exists: false},
	{name: "file1", isDir: false, exists: true},
	{name: "dir1", isDir: true, exists: true},
}

func TestIsDir(t *testing.T) {
	testIsDirFileCases(t, "IsDir", true, IsDir)
}

func TestIsFile(t *testing.T) {
	testIsDirFileCases(t, "IsFile", false, IsFile)
}

func TestCheckDir(t *testing.T) {
	testIsDirFileCases(t, "CheckDir", true, func(s string) bool { return CheckDir(s) == nil })
}

func TestCheckFile(t *testing.T) {
	testIsDirFileCases(t, "CheckFile", false, func(s string) bool { return CheckFile(s) == nil })
}

func testIsDirFileCases(t *testing.T, name string, checksDir bool, toTest func(string) bool) {
	for _, c := range isDirFileCases {
		if !c.exists {
			assert.False(t, toTest(path.Join(testDir, c.name)), "%s returns okay for a non-existent file", name)
		} else {
			assert.Equal(t, checksDir == c.isDir, toTest(path.Join(testDir, c.name)), "%s returns incorrect for a %s", name, dirFileName(c.isDir))
		}
	}
}

func dirFileName(isDir bool) string {
	if isDir {
		return "dir"
	}
	return "file"
}
