package routes

import (
	"git.peek1e.eu/peek1e/file-drop/server/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.RequestURI, "/dl/", "", -1)
	name, path, err := models.GetFileByID(id)
	if err != nil {
		log.Printf("Couldn't retrieve Database Entry %s", err)
		writeHeader(w, http.StatusNotFound)
		http.NotFound(w, r)
		return
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Couldn't retrieve File %s", err)
		writeHeader(w, http.StatusInternalServerError)
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
		log.Printf("Couldn't retrieve Database Entry %s", err)
		writeHeader(w, http.StatusNotFound)
		return
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Couldn't retrieve File %s", err)
		writeHeader(w, http.StatusInternalServerError)
		return
	}
	cType := http.DetectContentType(file)
	w.Header().Set("Content-Type", cType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+name+"\"")
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

func writeHeader(w http.ResponseWriter, statusCode int) {
	code := strconv.Itoa(statusCode)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(code))
}
