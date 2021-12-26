package main

import (
	"git.peek1e.eu/peek1e/file-drop/server/db"
	_ "git.peek1e.eu/peek1e/file-drop/server/db"
)

func init() {
	db.GetInstance()
}

func main() {

}
