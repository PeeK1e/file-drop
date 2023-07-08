package models

import (
	"log"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
	"gitlab.com/PeeK1e/file-drop/pkg/util"
)

func RemoveExpiredFiles() {
	statement := `DELETE FROM "file" WHERE "expirationtime" < NOW()::TIMESTAMP RETURNING "path"`

	r, err := db.RunTransactionWithResult([]byte(statement))
	if err != nil {
		log.Printf("WARN: Could not delete entries %s", err)
	}

	/*	TODO: cleaner way of deleting file, on failed delete should not remove entry from DB
	 *	Create A Log Entry in the database with failed file deletions?
	 */
	for r.Next() {
		path := ""
		r.Scan(&path)
		util.DeleteFile(path)
	}

	db.CloseRows(r)
}
