package helpers

import (
	"bufio"
	"os"
	"path/filepath"

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
		if err := os.WriteFile(destination, data, 0644); err != nil {
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

func MergeToFile(lines []string, file string) error {
	/*
		All lines that don't already exist in the file will be added to it
	*/
	readFile, err := os.Open(file)
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	readFile.Close()

	var c string
	for _, line := range RemoveDuplicates(append(fileLines, lines...)) {
		c += line + "\n"
	}

	return WriteFile(c, file)
}
