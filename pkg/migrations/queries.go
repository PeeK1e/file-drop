package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

const (
	DB_CLEAN int = 0
	DB_DIRTY int = 1
)

// This checks if the DB is initialized
func queryVersion() (int, error) {
	q := `SELECT value_int FROM meta WHERE alias LIKE 'migration_version'`

	// if there is an error it is likely that the DB is not initialized
	// assuming ther is no error here we return 0 and nil
	// --> migration version 0 initializes the DB
	r, err := db.GetInstance().Query(q)
	if err != nil {
		return 0, nil
	}

	// defer connection close
	defer func(r *sql.Rows) {
		if r == nil {
			fmt.Printf("WARN: Skipping close since Rows is nil %s", err)
		}
		if err := r.Close(); err != nil {
			fmt.Printf("WARN: error on rows close %s", err)
		}
	}(r)

	if !r.Next() {
		return 0, fmt.Errorf("no results in query")
	}

	var i int
	if err != r.Scan(&i) {
		return 0, fmt.Errorf("results not readable")
	}

	return i, nil
}

func queryIsDbDirty() bool {
	q := `SELECT value_int 
			FROM meta 
			WHERE alias LIKE 'dirty_migration'`

	// if there is an error it is likely that the DB is not initialized
	// assuming ther is no error here we return false to literally indicate a clean state
	r, err := db.GetInstance().Query(q)
	if err != nil {
		log.Printf("WARN: Couldn't fetch row %s, assuming DB is empty", err)
		return false
	}

	// defer connection close
	defer func(r *sql.Rows) {
		if r == nil {
			fmt.Printf("WARN: Skipping close since Rows is nil %s", err)
		}
		if err := r.Close(); err != nil {
			fmt.Printf("WARN: error on rows close %s", err)
		}
	}(r)

	// fetch next row
	if !r.Next() {
		log.Printf("WARN: empty result, assuming DB is empty")
		return false
	}

	var i int
	r.Scan(&i)

	return i >= 1
}

// Sets the `dirty_migration` flag.
// 0 equals clean state.
// 1 euals dirty state.
func setDirtyFlag(value int) {
	log.Printf("INFO: Setting Database State to %d", value)
	q := `UPDATE meta 
			SET value_int = $1 
			WHERE alias LIKE 'dirty_migration'`

	err := db.RunTransaction([]byte(q), value)
	if err != nil {
		log.Fatalf("ERR: Could not set dirty flag %s", err)
	}
}

// Sets the `migration_version` column int value
func setMigrationVersion(value int) {
	log.Printf("INFO: Setting Database Migration Level to %d", value)
	q := `UPDATE meta 
			SET value_int = $1 
			WHERE alias LIKE 'migration_version'`

	err := db.RunTransaction([]byte(q), value)
	if err != nil {
		log.Fatalf("ERR: Could not set migration version %s", err)
	}
}
