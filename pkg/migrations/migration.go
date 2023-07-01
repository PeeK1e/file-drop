package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

func Run() {
	db := db.GetInstance()

	if !checkDbInitState() {
		log.Fatal("ERR: DB Migration Failed")
		return
	}

	// Check if last migrations were done correctly
	if queryIsDbDirty(db) {
		log.Fatal("ERR: Database is dirty, aborting")
	}
	setDirtyFlag(db, 1)

	version, err := queryVersion(db)
	if err != nil {
		log.Fatal("ERR: no version detected")
		return
	}

	dirs, err := getDirs("./upgrade-db")
	if err != nil {
		log.Printf("ERR: Opening directory %s", err)
		return
	}

	runScripts(version, dirs, db)

	setDirtyFlag(db, 0)
}

func runScripts(version int, dirs []string, db *sql.DB) {
	for i, v := range dirs {
		// skip dir if less than migration version
		if i < version {
			continue
		}

		// read files and run transaction
		dir := fmt.Sprintf("./upgrade-db/%s", v)

		files, err := getFiles(dir)
		if err != nil {
			log.Printf("ERR: Opening file %s", err)
			return
		}

		for _, f := range files {
			filepath := fmt.Sprintf("%s/%s", dir, f)
			q, err := ioutil.ReadFile(filepath)
			if err != nil {
				log.Printf("Can't open migration file")
			}

			err = runTransaction(q, db)
			if err != nil {
				log.Printf("WARN: migration no. %d, named %s failed, reason %s", version, f, err)
			}
		}
		version++
		setMigrationVersion(db, version)
	}
}
