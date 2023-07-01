package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileHelpers(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "offlinerss-files")
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")

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

		lines, err = ReadFileToLines(file2)
		require.NoError(t, err)
		assert.Len(t, lines, 2)
	})

	t.Run("merge to file", func(t *testing.T) {
		err := MergeToFile([]string{
			"test1337",
			"test1337", // duplicate that should be removed
		}, file1)
		require.NoError(t, err)

		lines, err := ReadFileToLines(file1)
		require.NoError(t, err)
		assert.Len(t, lines, 3)
	})
}
