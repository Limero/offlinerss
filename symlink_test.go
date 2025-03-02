package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/limero/offlinerss/client"
	"github.com/limero/offlinerss/client/newsboat"
	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CaptureStdout(f func()) string {
	originalStdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = originalStdout
	}()

	f()

	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func Test_symlinkClientPaths(t *testing.T) {
	log.DebugEnabled = true
	sourceDir := t.TempDir()
	dataPath := domain.DataPath(sourceDir)
	sourceFile := "source-file"
	sourcePath := dataPath.GetFile(sourceFile)
	helpers.WriteFile("hello", sourcePath)

	setupClients := func(targetPaths ...string) []domain.Client {
		client := client.Client{
			DataPath: dataPath,
			Files: domain.ClientFiles{
				{
					FileName:    sourceFile,
					TargetPaths: targetPaths,
				},
			},
		}
		return []domain.Client{
			newsboat.Newsboat{client},
		}
	}

	t.Run("should link if link doesn't exist", func(t *testing.T) {
		targetPath := filepath.Join(t.TempDir(), "new-link")
		clients := setupClients(targetPath)

		output := CaptureStdout(func() {
			require.NoError(t, symlinkClientPaths(clients))
		})
		assert.Contains(t, output, "Symlinking")
		dest, err := os.Readlink(targetPath)
		require.NoError(t, err)
		assert.Equal(t, sourcePath, dest)
	})

	t.Run("should create subdirs and link if link and dirs doesn't exist", func(t *testing.T) {
		targetPath := filepath.Join(t.TempDir(), "sub/dir/new-link")
		clients := setupClients(targetPath)

		output := CaptureStdout(func() {
			require.NoError(t, symlinkClientPaths(clients))
		})
		assert.Contains(t, output, "Symlinking")
		dest, err := os.Readlink(targetPath)
		require.NoError(t, err)
		assert.Equal(t, sourcePath, dest)
	})

	t.Run("should link multiple if links doesn't exist", func(t *testing.T) {
		targetPath1 := filepath.Join(t.TempDir(), "new-link1")
		targetPath2 := filepath.Join(t.TempDir(), "new-link2")
		clients := setupClients([]string{targetPath1, targetPath2}...)

		output := CaptureStdout(func() {
			require.NoError(t, symlinkClientPaths(clients))
		})
		assert.Contains(t, output, "Symlinking")
		assert.Equal(t, 2, strings.Count(output, "Symlinking"))
		dest1, err := os.Readlink(targetPath1)
		require.NoError(t, err)
		assert.Equal(t, sourcePath, dest1)
		dest2, err := os.Readlink(targetPath1)
		require.NoError(t, err)
		assert.Equal(t, sourcePath, dest2)
	})

	t.Run("should skip if link exists and is correct", func(t *testing.T) {
		targetPath := filepath.Join(t.TempDir(), "existing-link")
		require.NoError(t, os.Symlink(sourcePath, targetPath))
		clients := setupClients(targetPath)

		output := CaptureStdout(func() {
			require.NoError(t, symlinkClientPaths(clients))
		})
		assert.Contains(t, output, "is already correct")
		dest, err := os.Readlink(targetPath)
		require.NoError(t, err)
		assert.Equal(t, sourcePath, dest)
	})

	t.Run("should unlink and link if link exists but is incorrect", func(t *testing.T) {
		incorrectSourcePath := dataPath.GetFile("incorrect")
		helpers.WriteFile("incorrect", incorrectSourcePath)

		targetPath := filepath.Join(t.TempDir(), "existing-incorrect-link")
		require.NoError(t, os.Symlink(incorrectSourcePath, targetPath))
		clients := setupClients(targetPath)

		output := CaptureStdout(func() {
			require.NoError(t, symlinkClientPaths(clients))
		})
		assert.Contains(t, output, "Removing incorrect symlink")
		assert.Contains(t, output, "Symlinking")
		dest, err := os.Readlink(targetPath)
		require.NoError(t, err)
		assert.Equal(t, sourcePath, dest)
	})

	t.Run("should rename and link if path exists and is file", func(t *testing.T) {
		targetPath := filepath.Join(t.TempDir(), "existing-file")
		helpers.WriteFile("different", targetPath)
		clients := setupClients(targetPath)

		output := CaptureStdout(func() {
			require.NoError(t, symlinkClientPaths(clients))
		})
		assert.Contains(t, output, "Non-symlink found at target")
		assert.Contains(t, output, "Symlinking")
		assert.True(t, helpers.FileExists(targetPath+".bak"))
		dest, err := os.Readlink(targetPath)
		require.NoError(t, err)
		assert.Equal(t, sourcePath, dest)
	})
}
