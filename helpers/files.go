package helpers

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/limero/offlinerss/log"
)

func CopyFile(source string, destinations ...string) error {
	/*
		This will copy a file from source to all destinations
	*/
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	for _, destination := range destinations {
		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
			return err
		}

		// Make a copy of the file
		if err := os.WriteFile(destination, data, 0o644); err != nil {
			return err
		}
		log.Debug("Copied file %s to %s", source, destination)
	}
	return nil
}

func WriteFile(content string, destinations ...string) error {
	/*
		Create file at destinations with content. If file exists, it will be overwritten.
	*/
	for _, destination := range destinations {
		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
			return err
		}

		file, err := os.Create(destination)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.WriteString(content); err != nil {
			return err
		}

		log.Debug("Wrote file: %s", destination)
	}

	return nil
}

func ReadFileToLines(file string) ([]string, error) {
	var fileLines []string
	readFile, err := os.Open(file)
	if err != nil {
		return fileLines, err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	readFile.Close()

	return fileLines, nil
}

func MergeToFile(lines []string, file string, sortFunc func(s1, s2 string) bool) error {
	/*
		All lines that don't already exist in the file will be added to it
		If original file doesn't exist, it will be created with the new lines
	*/
	fileLines, err := ReadFileToLines(file)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	lines = append(fileLines, lines...)

	if sortFunc != nil {
		sort.Slice(lines, func(i, j int) bool {
			return sortFunc(lines[i], lines[j])
		})
	}

	var c string
	for _, line := range RemoveDuplicates(lines) {
		c += line + "\n"
	}

	return WriteFile(c, file)
}
