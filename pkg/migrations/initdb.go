package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

func initDB(db *sql.DB) error {
	path := fmt.Sprintf("./upgrade-db/%d", 0)
	log.Print("INFO: Trying to initialize DB")
	f, err := getFiles(path)
	if err != nil {
		return err
	}
	for _, v := range f {
		log.Printf("INFO: Running Init Script %s", v)

		filepath := fmt.Sprintf("%s/%s", path, v)
		q, err := ioutil.ReadFile(filepath)
		if err != nil {
			return err
		}

		err = runTransaction(q, db)
		if err != nil {
			return err
		}
	}
	return nil
}
