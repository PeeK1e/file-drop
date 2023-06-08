package util

import (
	"log"
	"os"
)

func DeleteFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		log.Printf("Could not delete file at %s :: %s", path, err)
		return false
	}
	log.Printf("Deleted File %s", path)
	return true
}
