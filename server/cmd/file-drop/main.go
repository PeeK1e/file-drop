package main

import (
	"flag"
	"git.peek1e.eu/peek1e/file_drop/db"
	_ "git.peek1e.eu/peek1e/file_drop/db"
	_ "github.com/lib/pq"
	"log"
)

type settings struct {
	Port         string
	Address      string
	ServeAddress string
}

var setting settings

func init() {
	flag.StringVar(&setting.Port, "p", "8080", "Sets the listening port of the Webserver")
	flag.StringVar(&setting.Address, "add", "", "Sets a specific address to listen on e.g. only localhost")
}

func init() {
	err := db.GetInstance().Ping()
	if err != nil {
		log.Panic("DB not available")
	}
	log.Print("Connected to DB")
}

func main() {

}
