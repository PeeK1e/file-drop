package migrations

import (
	"fmt"
	"log"
	"os"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

func Run() {
	// Check if last migrations were done correctly
	// abort when not in clean state
	if queryIsDbDirty() {
		log.Fatal("ERR: Database is dirty, aborting")
	}

	version, err := queryVersion()
	if err != nil {
		log.Fatalf("ERR: no version detected %s", err)
	}

	dirs, err := getDirs("./upgrade-db")
	if err != nil {
		log.Fatalf("ERR: Opening directory %s", err)
	}

	runScripts(version, dirs)
}

func runScripts(version int, dirs []string) {
	for i, v := range dirs {
		// skip dir if less than migration version
		if i < version {
			continue
		}

		// read files and run transaction
		dir := fmt.Sprintf("./upgrade-db/%s", v)

		files, err := getFiles(dir)
		if err != nil {
			log.Fatalf("ERR: Opening file %s", err)
			return
		}

		// version 0 equals an empty database
		if version >= 1 {
			setDirtyFlag(DB_DIRTY)
		}

		for _, f := range files {
			filepath := fmt.Sprintf("%s/%s", dir, f)
			q, err := os.ReadFile(filepath)
			if err != nil {
				log.Fatalf("ERR: Can't open migration file")
			}
			log.Printf("INFO: Migration Level: %d, Running Script %s", version, f)
			err = db.RunTransaction(q)
			if err != nil {
				log.Fatalf("ERR: migration no. %d, named %s failed, reason %s", version, f, err)
			}
		}

		version++
		setMigrationVersion(version)
		setDirtyFlag(DB_CLEAN)
	}
}
