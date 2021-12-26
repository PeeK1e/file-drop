package main

import (
	"flag"
	"fmt"
	"git.peek1e.eu/peek1e/file_drop/db"
	"git.peek1e.eu/peek1e/file_drop/routes"
	_ "git.peek1e.eu/peek1e/file_drop/routes"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Settings struct {
	Port         string
	Address      string
	ServeAddress string
	//DataDir      string
}

var setting Settings

func init() {
	fmt.Println("STARTING UP")
	flag.StringVar(&setting.Port, "p", "8080", "Sets the listening port of the Webserver")
	flag.StringVar(&setting.Address, "addr", "", "Sets a specific address to listen on e.g. only localhost")
	//flag.StringVar(&setting.DataDir, "d", "./storage", "Sets the Data Directory Path where files are uploaded to")
	setting.ServeAddress = setting.Address + ":" + setting.Port

	err := os.Mkdir("./storage", 0764)
	if err != nil {
		log.Printf("Cant create Data Dir %s", err)
	}
	rand.Seed(time.Now().UTC().UnixNano())
}

func init() {
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

func main() {
	fs := http.FileServer(http.Dir("./upload_html"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.HandleFunc("/dl/", routes.DownloadFile)
	http.HandleFunc("/pv/", routes.PreviewFile)
	http.HandleFunc("/upload", routes.UploadFile)

	http.ListenAndServe(setting.ServeAddress, nil)
}
