package models

import (
	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

func SaveFile(fileID string, fileName string, filePath string, encrypted bool, sha string) error {
	statement := `INSERT INTO "file" ("keyid", "filename", "path", "is_encrypted", "secret_sha") VALUES ($1,$2,$3,$4,$5)`
	return db.RunTransaction([]byte(statement), fileID, fileName, filePath, encrypted, sha)
}
