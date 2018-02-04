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

func TestIsDir(t *testing.T) {
	assert.False(t, IsDir(path.Join(testDir, "non-existent")), "IsDir returns true for a non-existent file")
	assert.False(t, IsDir(path.Join(testDir, "file1")), "IsDir returns true for a file")
	assert.True(t, IsDir(path.Join(testDir, "dir1")), "IsDir returns false for a dir")
}

func TestIsFile(t *testing.T) {
	assert.False(t, IsFile(path.Join(testDir, "non-existent")), "IsDir returns true for a non-existent file")
	assert.True(t, IsFile(path.Join(testDir, "file1")), "IsDir returns false for a file")
	assert.False(t, IsFile(path.Join(testDir, "dir1")), "IsDir returns true for a dir")
}

func TestCheckDir(t *testing.T) {
	assert.NotNil(t, CheckDir(path.Join(testDir, "non-existent")), "CheckDir returned no error for a non-existent file")
	assert.NotNil(t, CheckDir(path.Join(testDir, "file1")), "CheckDir returned no error for a file")
	assert.Nil(t, CheckDir(path.Join(testDir, "dir1")), "CheckDir returned an error for a dir")
}

func TestCheckFile(t *testing.T) {
	assert.NotNil(t, CheckFile(path.Join(testDir, "non-existent")), "CheckFile returned no error for a non-existent file")
	assert.Nil(t, CheckFile(path.Join(testDir, "file1")), "CheckFile returned an error for a file")
	assert.NotNil(t, CheckFile(path.Join(testDir, "dir1")), "CheckFile returned no error for a dir")
}
