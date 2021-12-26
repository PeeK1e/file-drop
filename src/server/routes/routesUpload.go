package routes

import (
	"encoding/json"
	"git.peek1e.eu/peek1e/file_drop/models"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type uploadResponse struct {
	Reason string `json:"reason"`
	Ok     bool   `json:"Ok"`
	FileID string `json:"fileID"`
}

const (
	letters = string("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100000000)
	tFile, fileHeader, err := r.FormFile("file")
	defer tFile.Close()
	if err != nil {
		log.Printf("Could not parse File! %s", err)
		sendResponse(w, uploadResponse{
			Reason: "Server Error",
			Ok:     false,
			FileID: "",
		})
		return
	}

	filePath, dirPathChild := getRandomPathName()
	id := createRandomId(10)

	err = os.MkdirAll(dirPathChild, 0764)
	if err != nil {
		log.Printf("Error Creating Directory %s", err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	_, err = io.Copy(file, tFile)
	if err != nil {
		log.Printf("Error Writing File %s", err)
		sendResponse(w, uploadResponse{
			Reason: "Server Error",
			Ok:     false,
			FileID: "",
		})
		return
	}

	err = models.SaveFile(id, fileHeader.Filename, filePath)
	if err != nil {
		log.Printf("Error Saving File to DB %s,", err)
		sendResponse(w, uploadResponse{
			Reason: "Database Error",
			Ok:     false,
			FileID: "",
		})
		_ = os.Remove(filePath)
		return
	}

	log.Printf("File %s", fileHeader.Filename)
	sendResponse(w, uploadResponse{
		Reason: "File OK",
		Ok:     true,
		FileID: id,
	})
}

func sendResponse(w http.ResponseWriter, uR uploadResponse) {
	jsString, err := json.Marshal(uR)
	if err != nil {
		log.Printf("Could not marshal struct %s", err)
	}
	_, err = io.WriteString(w, string(jsString))
	if err != nil {
		log.Printf("Could not send response %s", err)
	}
}

func getRandomPathName() (string, string) {
	var pathParts [3]string

	for i := 0; i < 2; i++ {
		pathParts[0] += string(letters[rand.Intn(len(letters))])
	}
	for i := 0; i < 2; i++ {
		pathParts[1] += string(letters[rand.Intn(len(letters))])
	}
	for i := 0; i < 10; i++ {
		pathParts[2] += string(letters[rand.Intn(len(letters))])
	}

	dirPathChild := "./storage/" + pathParts[0] + "/" + pathParts[1]
	filePath := "./storage/" + pathParts[0] + "/" + pathParts[1] + "/" + pathParts[2]

	ok, err := models.IsPathOk(filePath)
	if !ok {
		log.Printf("Already exists or an error occured: %s", err)
		return getRandomPathName()
	}

	return filePath, dirPathChild
}

func createRandomId(length int) string {
	str := ""
	for i := 0; i < length; i++ {
		str += string(letters[rand.Intn(len(letters))])
	}

	ok, err := models.IsFileIdOk(str)
	if !ok {
		log.Printf("Already exists or an error occured: %s", err)
		return createRandomId(length)
	}

	return str
}
