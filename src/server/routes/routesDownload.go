package routes

import (
	"git.peek1e.eu/peek1e/file_drop/server/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.RequestURI, "/dl/", "", -1)
	name, path, err := models.GetFileByID(id)
	if err != nil {
		log.Printf("Couldn't retrieve File %s", err)
		_, _ = io.WriteString(w, "HTTP 404")
		return
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Couldn't retrieve File %s", err)
		_, _ = io.WriteString(w, "HTTP 500")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(file)
}

func PreviewFile(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.RequestURI, "/pv/", "", -1)
	name, path, err := models.GetFileByID(id)
	if err != nil {
		log.Printf("Couldn't retrieve File %s", err)
		_, _ = io.WriteString(w, "HTTP 404")
		return
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Couldn't retrieve File %s", err)
		_, _ = io.WriteString(w, "HTTP 500")
		return
	}
	cType := http.DetectContentType(file)
	w.Header().Set("Content-Type", cType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+name+"\"")
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}
