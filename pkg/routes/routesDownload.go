package routes

import (
	"log"
	"net/http"
	"os"
	"strings"

	"gitlab.com/PeeK1e/file-drop/pkg/models"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.RequestURI, "/pv/", "", -1)

	log.Printf("Client %s requested FileID %s", r.RemoteAddr, id)

	_, path, err := models.GetFileByID(id)
	if err != nil {
		log.Printf("Couldn't retrieve Database Entry %s", err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0640)
	if err != nil {
		log.Printf("Couldn't retrieve File %s", err)
		http.Error(w, "File not found on server", http.StatusInternalServerError)
		return
	}

	fileStat, err := file.Stat()
	if err != nil {
		log.Printf("File Corrupted %s", err)
		http.Error(w, "File Corrupted", http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, path, fileStat.ModTime(), file)
}
