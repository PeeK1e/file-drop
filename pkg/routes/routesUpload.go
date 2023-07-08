package routes

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"gitlab.com/PeeK1e/file-drop/pkg/models"
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
	r.ParseMultipartForm(1000000000)
	tFile, fileHeader, err := r.FormFile("file")

	log.Printf("New File Upload %s", fileHeader.Filename)

	//dont care about the error, function will fail if upload wasn't successfull
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

	if err = os.MkdirAll(dirPathChild, 0760); err != nil {
		log.Printf("Error Creating Directory %s", err)
	}

	log.Printf("Saving File in %s", filePath)

	file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0660)
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

	filename := strings.Split(fileHeader.Filename, ".")
	fileExt := filename[len(filename)-1]
	id := ""
	if filename[0] != fileExt {
		id = createRandomId(10, fileExt)
	} else {
		id = createRandomId(10, "")
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

	log.Printf("Saved File %s, in %s", fileHeader.Filename, filePath)
}

func sendResponse(w http.ResponseWriter, uR uploadResponse) {
	jsString, err := json.Marshal(uR)
	if err != nil {
		log.Printf("Could not marshal struct %s", err)
		jsString = []byte("{}")
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

	if !models.IsPathOk(filePath) {
		log.Printf("WARN: Path not ok, regenerating...")
		return getRandomPathName()
	}

	return filePath, dirPathChild
}

func createRandomId(length int, fileExt string) string {
	str := ""
	for i := 0; i < length; i++ {
		str += string(letters[rand.Intn(len(letters))])
	}
	if fileExt != "" {
		str += "." + fileExt
	}

	if !models.IsFileIdOk(str) {
		log.Printf("WARN: id not ok, regenerating...")
		return createRandomId(length, fileExt)
	}

	return str
}
