package models

import (
	"errors"
	"fmt"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

func GetFileByID(id string) (name string, path string, err error) {
	instance := db.GetInstance()

	statement := `SELECT "filename", "path" FROM file WHERE "keyid" LIKE $1`

	result, err := instance.Query(statement, id)
	if err != nil {
		return "", "", fmt.Errorf("couldn't exect query %s", err)
	}

	defer result.Close()
	for result.Next() {
		filename, filepath := "", ""
		err = result.Scan(&filename, &filepath)
		return filename, filepath, err
	}
	return "", "", errors.New("no such result")
}
