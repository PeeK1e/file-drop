package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/PeeK1e/file-drop/pkg/config"
	"gitlab.com/PeeK1e/file-drop/pkg/db"
	"gitlab.com/PeeK1e/file-drop/pkg/routes"
)

const (
	storagePath = string("/app/storage")
)

func main() {

	log.Println("INFO: Starting Server")

	// creating directory if non existent
	err := os.Mkdir(storagePath, 0764)
	if err != nil {
		log.Printf("WARN: Not creating directory, %s", err)
	}

	c := config.NewConfig()

	// Initialise the Database
	for !db.NewDB(*c.DbSettings) {
		return
	}

	http.HandleFunc("/pv/", routes.DownloadFile)
	http.HandleFunc("/upload", routes.UploadFile)
	http.HandleFunc("/healthz", routes.Healthz)

	go startHttpServer(*c.HttpServer.ListenAddress)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	<-sigChannel
	log.Print("INFO: Caught shutdown signal. Terminating.")
}

func startHttpServer(addr string) {
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("ERR: HTTP Server error, %s", err)
		os.Exit(500)
	}
}
