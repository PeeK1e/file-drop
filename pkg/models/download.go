package models

import (
	"fmt"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

func FileByID(id string) (name string, path string, err error) {
	statement := `SELECT "filename", "path" FROM file WHERE "keyid" LIKE $1`

	r, err := db.QueryWithResult([]byte(statement), id)
	if err != nil {
		return "", "", err
	}
	defer db.CloseRows(r)

	for r.Next() {
		filename, filepath := "", ""
		err = r.Scan(&filename, &filepath)
		return filename, filepath, err
	}

	return "", "", fmt.Errorf("no such result")
}

func EncryptionDetails(id string) (string, bool, error) {
	statement := `SELECT "secret_sha", "is_encrypted" FROM file WHERE "keyid" LIKE $1`

	r, err := db.QueryWithResult([]byte(statement), id)
	if err != nil {
		return "", false, err
	}
	defer db.CloseRows(r)

	sha, enc := "", false
	for r.Next() {
		err = r.Scan(&sha, &enc)
		if err != nil {
			return "", false, err
		}
	}

	return sha, enc, nil
}
