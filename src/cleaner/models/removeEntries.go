package models

import (
	"database/sql"
	"git.peek1e.eu/peek1e/file-drop/cleaner/db"
	"git.peek1e.eu/peek1e/file-drop/cleaner/util"
	"log"
)

func RemoveExpiredFiles() {
	instance := db.GetInstance()

	t, err := instance.Begin()
	if err != nil {
		log.Printf("Error creating transaction %s", err)
		return
	}

	statemnt := `DELETE FROM "file" WHERE "expirationtime" < NOW()::TIMESTAMP RETURNING "path"`

	rows, err := t.Query(statemnt)
	defer rows.Close()
	if err != nil {
		log.Printf("Could not execute statement %s", err)
		rollback(t)
		return
	}

	for rows.Next() {
		path := ""
		rows.Scan(&path)
		util.DeleteFile(path)
	}

	err = t.Commit()
	if err != nil {
		log.Printf("Could not commit transaction %s", err)
		rollback(t)
		return
	}
}

func rollback(t *sql.Tx) {
	err := t.Rollback()
	if err != nil {
		log.Printf("Rollback failed %s", err)
	}
}
