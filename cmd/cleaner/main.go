package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/PeeK1e/file-drop/pkg/config"
	"gitlab.com/PeeK1e/file-drop/pkg/db"
	"gitlab.com/PeeK1e/file-drop/pkg/models"
	"gitlab.com/PeeK1e/file-drop/pkg/routes"
)

func main() {

	c := config.NewConfig()

	// Initialise the Database
	for !db.NewDB(*c.DbSettings) {
		return
	}

	ticker := time.Ticker{C: time.Tick(1 * time.Minute)}

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// handle health checks
	http.HandleFunc("/healthz", routes.Healthz)
	go startHttpServer(*c.HttpServer.ListenAddress)

	for {
		select {
		case <-ticker.C:
			log.Printf("INFO: Starting cleaner...")
			models.RemoveExpiredFiles()
			log.Printf("INFO: Done.")
		case <-sigChannel:
			log.Print("INFO: Caught shutdown signal. Terminating.")
			return
		}
	}
}

func startHttpServer(addr string) {
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("ERR: HTTP Server error, %s", err)
		os.Exit(500)
	}
}
