package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/PeeK1e/file-drop/pkg/config"
	"gitlab.com/PeeK1e/file-drop/pkg/db"
	"gitlab.com/PeeK1e/file-drop/pkg/models"
)

func main() {

	c := config.NewConfig()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// Initialise the Database
	for !db.NewDB(*c.DbSettings) {
		return
	}

	ticker := time.Ticker{C: time.Tick(10 * time.Minute)}

	select {
	case <-ticker.C:
		models.RemoveExpiredFiles()
	case <-sigChannel:
		log.Print("Caught shutdown signal. Terminating.")
		return
	}
}
