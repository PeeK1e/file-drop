package migrations

import (
	"os"
	"sort"
)

// reading all files from a directory
//
// returns a list of files or error if reading fails
func getFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}

	sort.Strings(filenames)

	return filenames, nil
}

// reading all directories from given path
//
// returns a list of directories or error if reading fails
func getDirs(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirnames []string
	for _, file := range files {
		if file.IsDir() {
			dirnames = append(dirnames, file.Name())
		}
	}

	sort.Strings(dirnames)

	return dirnames, nil
}
