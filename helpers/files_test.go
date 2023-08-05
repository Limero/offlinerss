package helpers

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileHelpers(t *testing.T) {
	file1 := filepath.Join(t.TempDir(), "file1.txt")
	file2 := filepath.Join(t.TempDir(), "file2.txt")

	t.Run("write file", func(t *testing.T) {
		err := WriteFile("test123\ntest321", file1)
		require.NoError(t, err)
	})

	t.Run("copy file", func(t *testing.T) {
		err := CopyFile(file1, file2)
		require.NoError(t, err)
	})

	t.Run("read files", func(t *testing.T) {
		lines, err := ReadFileToLines(file1)
		require.NoError(t, err)
		assert.Len(t, lines, 2)
		assert.Equal(t, []string{"test123", "test321"}, lines)

		lines, err = ReadFileToLines(file2)
		require.NoError(t, err)
		assert.Len(t, lines, 2)
		assert.Equal(t, []string{"test123", "test321"}, lines)
	})

	t.Run("merge to file", func(t *testing.T) {
		err := MergeToFile([]string{
			"test1337",
			"test1337", // duplicate that should be removed
		}, file1, nil)
		require.NoError(t, err)

		lines, err := ReadFileToLines(file1)
		require.NoError(t, err)
		assert.Len(t, lines, 3)
		assert.Equal(t, []string{"test123", "test321", "test1337"}, lines)
	})

	t.Run("merge to file with sort", func(t *testing.T) {
		sortFunc := func(s1, s2 string) bool {
			return s1 < s2
		}

		err := MergeToFile([]string{
			"test322",
		}, file1, sortFunc)
		require.NoError(t, err)

		lines, err := ReadFileToLines(file1)
		require.NoError(t, err)
		assert.Len(t, lines, 4)
		assert.Equal(t, []string{"test123", "test1337", "test321", "test322"}, lines)
	})
}
