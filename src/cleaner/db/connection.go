package db

import (
	"cleaner/environment"
	"database/sql"
	"log"
	"sync"
)

var (
	db   *sql.DB
	lock *sync.Mutex
)

func init() {
	lock = new(sync.Mutex)
}

func GetInstance() *sql.DB {
	lock.Lock()
	defer lock.Unlock()

	if db == nil {
		connectionString := environment.DatabaseSettings.GetDatabaseString()

		err := error(nil)
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Panicf("Could not open Connection to Database %s", err)
		}
	}

	return db
}
