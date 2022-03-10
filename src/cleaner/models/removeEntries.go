package models

import (
	"cleaner/db"
	"cleaner/util"
	"database/sql"
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

	/*	TODO: cleaner way of deleting file, on failed delete should not remove entry from DB
	 *	Create A Log Entry in the database with failed file deletions?
	 */
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
