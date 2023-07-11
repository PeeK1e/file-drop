package routes

import (
	"crypto/sha512"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gitlab.com/PeeK1e/file-drop/pkg/models"
	"gitlab.com/PeeK1e/file-drop/pkg/util"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.RequestURI, "/pv/", "", -1)

	log.Printf("Client %s requested FileID %s", r.RemoteAddr, id)

	_, path, err := models.FileByID(id)
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

func DownloadEncryptedFile(w http.ResponseWriter, r *http.Request) {
	//?needs refactoring later
	//id := strings.Replace(r.RequestURI, "/enc/", "", -1)
	log.Printf("Client %s requested encrypted file", r.RemoteAddr)

	wsChallenge(w, r)
}

// opens a websocket and creates the challenge for the client to calculate
func wsChallenge(w http.ResponseWriter, r *http.Request) {
	// opens websocket, if fails returns HTTP error code
	// like required by the RFC
	sock := &socket{}
	if !sock.NewSocket(w, r) {
		return
	}

	defer sock.Close()

	var fileID string
	if ok := sock.ReadClientData(&fileID); !ok {
		log.Printf("ERR: WS read stream broken")
		return
	}

	log.Printf("DEBUG: file id: %s // %s", fileID, string(fileID))

	path, challenge, challengeSha, err := createChallengeKey(string(fileID))
	if err != nil {
		log.Printf("ERR: Could not retrieve file %s", err)
		return
	}

	if ok := sock.WriteMessage(challenge); !ok {
		log.Printf("ERR: WS write stream broken")
		return
	}

	var challengeResp string
	if ok := sock.ReadClientData(&challengeResp); !ok {
		log.Printf("ERR: WS read stream broken")
		return
	}

	if challengeResp != string(challengeSha) {
		log.Printf("WARN: Client response not matching with server challenge")
		return
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}

	if ok := sock.WriteBytes(fileData); !ok {
		log.Printf("ERR: WS write stream broken")
		return
	}
}

// retrieves the sha and path for an file id
func shaForID(id string) (string, string) {
	sha, isEnc, err := models.EncryptionDetails(id)
	if !isEnc || err != nil {
		log.Printf("ERR: file doesn't look encrypted, aborting...")
		return "", ""
	}

	_, path, err := models.FileByID(id)
	if err != nil {
		log.Printf("Couldn't retrieve Database Entry %s", err)
		//http.Error(w, "File not found", http.StatusNotFound)
		return "", ""
	}

	return sha, path
}

// this creates the challenge key and retrieves informationm about the file
//
// returns the file path, the challenge, challenge SHA512
// returns an error if the file ID couldn't be found
func createChallengeKey(fileID string) (string, string, []byte, error) {
	pwSha, filePath := shaForID(string(fileID))
	if pwSha == "" {
		return "", "", []byte(""), fmt.Errorf("ERR: Could not open file")
	}

	//create challenge
	challenge := util.GetRandomString(12)
	sha_512 := sha512.New()
	sha_512.Write([]byte(pwSha + challenge))
	challengeSha := sha_512.Sum(nil)

	return filePath, challenge, challengeSha, nil
}
