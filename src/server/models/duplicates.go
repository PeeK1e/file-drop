package models

import (
	"errors"
	"fmt"
	"server/db"
)

func IsPathOk(path string) (bool, error) {
	instance := db.GetInstance()

	statement := `SELECT COUNT("path") AS COUNT FROM "file" WHERE "path" LIKE $1`

	result, err := instance.Query(statement, path)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Can't exec statement %s", err))
	}

	defer result.Close()
	for result.Next() {
		count := -1
		_ = result.Scan(&count)
		if count > 0 {
			return false, nil
		}
	}

	return true, nil
}

func IsFileIdOk(id string) (bool, error) {
	instance := db.GetInstance()

	statement := `SELECT COUNT("keyid") AS CNT FROM "file" WHERE "keyid" LIKE $1`

	result, err := instance.Query(statement, id)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Can't exec statement %s", err))
	}

	defer result.Close()
	for result.Next() {
		count := -1
		_ = result.Scan(&count)
		if count > 0 {
			return false, nil
		}
	}

	return true, nil
}
