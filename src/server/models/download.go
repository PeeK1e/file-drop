package models

import (
	"errors"
	"fmt"
	"git.peek1e.eu/peek1e/file_drop/server/db"
)

func GetFileByID(id string) (name string, path string, err error) {
	instance := db.GetInstance()

	statement := `SELECT "filename", "path" FROM file WHERE "keyid" LIKE $1`

	result, err := instance.Query(statement, id)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("Couldn't exect query %s", err))
	}

	defer result.Close()
	for result.Next() {
		filename, filepath := "", ""
		err = result.Scan(&filename, &filepath)
		return filename, filepath, err
	}
	return "", "", errors.New("no such result")
}
