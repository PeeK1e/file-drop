package models

import (
	"errors"
	"fmt"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

func SaveFile(fileID string, fileName string, filePath string) error {
	instance := db.GetInstance()

	t, err := instance.Begin()
	if err != nil {
		return errors.New(fmt.Sprintf("Can't start transaction %s", err))
	}

	statement := `INSERT INTO "file" ("keyid", "filename", "path") VALUES ($1,$2,$3)`

	_, err = t.Exec(statement, fileID, fileName, filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Can't create statement %s", err))
	}

	err = t.Commit()
	if err != nil {
		return errors.New(fmt.Sprintf("Can't commit transaction %s", err))
	}

	return nil
}
