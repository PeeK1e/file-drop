package migrations

import (
	"io/ioutil"
	"sort"
)

func getFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
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

func getDirs(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
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
