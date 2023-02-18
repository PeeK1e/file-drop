package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"server/db"
	"server/environment"
	"server/routes"
	"time"

	_ "github.com/lib/pq"
)

const (
	storagePath = string("/app/storage")
)

func init() {
	log.Println("STARTING UP")

	// creating directory if non existent
	err := os.Mkdir(storagePath, 0764)
	if err != nil {
		log.Printf("Not creating directory, %s", err)
	}

	// Initialise the Database
	connectToDatabase()

	// Initialising random seed
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	http.HandleFunc("/pv/", routes.DownloadFile)
	http.HandleFunc("/upload", routes.UploadFile)

	http.ListenAndServe(environment.HttpServerOptions.String(), nil)
}

func connectToDatabase() {
	count := 0
	duration := 2
	retryLimit := 5
	for {
		err := db.GetInstance().Ping()
		if err != nil && count <= retryLimit {
			log.Printf("Connection failed with error { %s } sleeping for %d, retry %d/%d", err, duration, count, retryLimit)
			time.Sleep(time.Duration(duration * int(time.Second)))
			count++
		} else if count > 5 {
			log.Panicf("DB not available %s", err)
		} else {
			log.Print("Connected to DB")
			break
		}
	}
}
