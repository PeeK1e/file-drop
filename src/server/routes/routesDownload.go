package routes

import (
	"server/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	file, err1 := ioutil.ReadFile(path)
	f, err2 := os.OpenFile(path, os.O_RDONLY, 0640)
	stat, err3 := f.Stat()
	defer f.Close()
	if err1 != nil || err2 != nil || err3 != nil {
		log.Printf("Couldn't retrieve File %s :: %s :: %s", err1, err2, err3)
		writeHeader(w, http.StatusInternalServerError)
		return
	}

	cType := http.DetectContentType(file)
	w.Header().Set("Content-Type", cType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+name+"\"")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(file)
}

func writeHeader(w http.ResponseWriter, statusCode int) {
	code := strconv.Itoa(statusCode)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(code))
}
