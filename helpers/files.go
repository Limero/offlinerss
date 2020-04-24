package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CopyFile(source string, destinations ...string) error {
	/*
		This will copy a file from source to all destinations
	*/
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	for _, destination := range destinations {
		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
			return err
		}

		// Make a copy of the file
		if err := ioutil.WriteFile(destination, data, 0644); err != nil {
			return err
		}
		fmt.Printf("Copied file %s to %s\n", source, destination)
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

		fmt.Printf("Wrote file: %s\n", destination)
	}

	return nil
}
