package db

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"gitlab.com/PeeK1e/file-drop/pkg/config"
)

var (
	psqlConnectionString string
	db                   *sql.DB
	lock                 *sync.Mutex
)

// Create new DB connection
// returns true on successfull connection
func NewDB(c config.DatabaseSettings) bool {

	lock = &sync.Mutex{}

	psqlConnectionString = c.GetDatabaseString()

	// retry db connection n times
	// where n is defined by the configuration settings
	count := 0
	for nil == GetInstance() {
		time.Sleep(time.Second)
		if count >= *c.ConnectionRetryCount {
			log.Printf("Could not open databse connection to host %s on port %s", *c.HOSTNAME, *c.PORT)
			return false
		}
		count++
	}

	if err := GetInstance().Ping(); err != nil {
		log.Printf("db conn not alive reason: %s", err)
		return false
	}

	log.Printf("db connection established")

	return true
}

func GetInstance() *sql.DB {
	lock.Lock()
	defer lock.Unlock()

	if db == nil {

		err := error(nil)
		db, err = sql.Open("postgres", psqlConnectionString)
		if err != nil {
			log.Printf("Could not open Connection to Database %s", err)
		}
	}

	return db
}
