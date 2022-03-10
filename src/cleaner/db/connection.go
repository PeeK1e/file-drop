package db

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

var (
	db        *sql.DB
	dbSetting *dbSettings
	lock      *sync.Mutex
	path      string
)

type dbSettings struct {
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	UserName   string `json:"userName"`
	UserPasswd string `json:"userPasswd"`
	DBName     string `json:"DBName"`
	SSLMode    string `json:"SSLMode"`
}

func init() {
	flag.StringVar(&path, "db-conf", "./db_settings/dbSettings.json", "Config file location for the Database")
	flag.Parse()

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("Make sure the dbSettings.json file exists %s", err)
	}

	dbSetting = new(dbSettings)
	err = json.Unmarshal(file, dbSetting)
	if err != nil {
		log.Panicf("Please make sure your json file has the right format %s", err)
	}
	lock = new(sync.Mutex)
}

func GetInstance() *sql.DB {
	if db == nil {
		lock.Lock()
		defer lock.Unlock()

		connString := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbSetting.Hostname,
			dbSetting.Port,
			dbSetting.UserName,
			dbSetting.UserPasswd,
			dbSetting.DBName,
			dbSetting.SSLMode)

		err := error(nil)
		db, err = sql.Open("postgres", connString)
		if err != nil {
			log.Panicf("Could not open Connection to Database %s", err)
		}
	}

	return db
}
