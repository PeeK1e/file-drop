package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

// This checks if the DB is initialized
func checkDbInitState() bool {
	_, err := queryVersion(db.GetInstance())
	if err == nil {
		return true
	}

	log.Printf("WARN: Database likely not initialized %s", err)

	// Try to init DB
	err = initDB(db.GetInstance())
	if err != nil {
		log.Printf("ERR: Failed to init DB %s", err)
		return false
	}

	return true
}

func queryVersion(db *sql.DB) (int, error) {
	q := `SELECT value_int FROM meta WHERE alias LIKE 'migration_version'`

	// if there is an error it is likely that the DB is not initialized
	r, err := db.Query(q)
	//defer r.Close()
	if err != nil {
		return 0, err
	}

	if !r.Next() {
		return 0, fmt.Errorf("no results in query")
	}
	var i int
	r.Scan(&i)

	return i, nil
}

func queryIsDbDirty(db *sql.DB) bool {
	q := `SELECT value_int FROM meta WHERE alias LIKE 'dirty_migration'`
	r, err := db.Query(q)
	//defer r.Close()
	if err != nil {
		log.Printf("ERR: Couldn't fetch row %s", err)
		return true
	}

	if !r.Next() {
		log.Printf("WARN: empty result, assuming DB is empty")
		return true
	}

	var i int
	r.Scan(&i)

	return i >= 1
}

// Sets the `dirty_migration` flag
// 0 == clean
// 1 == dirty
func setDirtyFlag(db *sql.DB, value int) {
	q := `UPDATE meta 
			SET value_int = $1 
			WHERE alias LIKE 'dirty_migration'`

	_, err := db.Exec(q, value)
	if err != nil {
		log.Fatalf("ERR: Could not set dirty flag %s", err)
	}
}

// Sets the `migration_version`
func setMigrationVersion(db *sql.DB, value int) {
	q := `UPDATE meta 
			SET value_int = $1 
			WHERE alias LIKE 'migration_version'`

	_, err := db.Exec(q, value)
	if err != nil {
		log.Fatalf("ERR: Could not set migration version %s", err)
	}
}

func runTransaction(q []byte, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(q))
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
