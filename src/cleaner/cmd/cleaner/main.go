package main

import (
	"cleaner/db"
	"cleaner/models"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func init() {
	connectToDatabase()
}

func main() {
	ticker := time.Ticker{C: time.Tick(10 * time.Minute)}
	for {
		<-ticker.C
		models.RemoveExpiredFiles()
	}
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
