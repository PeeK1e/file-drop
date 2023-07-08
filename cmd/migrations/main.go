package main

import (
	"log"
	"os"

	"gitlab.com/PeeK1e/file-drop/pkg/config"
	"gitlab.com/PeeK1e/file-drop/pkg/db"
	"gitlab.com/PeeK1e/file-drop/pkg/migrations"
)

func main() {
	c := config.NewConfig()

	if !db.NewDB(*c.DbSettings) {
		log.Print("ERR: Could not open DB connection")
		os.Exit(1)
	}

	log.Printf("INFO: Running Database Migrations...")

	migrations.Run()

	log.Printf("INFO: Done! :)")
}
