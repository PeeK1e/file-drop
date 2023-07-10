package models

import (
	"log"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

func IsPathOk(path string) bool {
	statement := `SELECT COUNT("path") AS COUNT FROM "file" WHERE "path" LIKE $1`

	r, err := db.QueryWithResult([]byte(statement), path)
	defer db.CloseRows(r)
	if err != nil {
		log.Printf("WARN: Could not run query %s", err)
		return false
	}

	for r.Next() {
		count := -1
		_ = r.Scan(&count)
		if count > 0 {
			return false
		}
	}

	return true
}

func IsFileIdOk(id string) bool {
	statement := `SELECT COUNT("keyid") AS CNT FROM "file" WHERE "keyid" LIKE $1`

	r, err := db.QueryWithResult([]byte(statement), id)
	defer db.CloseRows(r)
	if err != nil {
		log.Printf("WARN: Could not run query %s", err)
		return false
	}

	for r.Next() {
		count := -1
		_ = r.Scan(&count)
		if count > 0 {
			return false
		}
	}

	return true
}
