package models

import (
	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

func SaveFile(fileID string, fileName string, filePath string) error {
	statement := `INSERT INTO "file" ("keyid", "filename", "path") VALUES ($1,$2,$3)`
	return db.RunTransaction([]byte(statement), fileID, fileName, filePath)
}
