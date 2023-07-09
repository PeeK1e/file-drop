package routes

import (
	"crypto/sha512"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gobwas/ws"
	"gitlab.com/PeeK1e/file-drop/pkg/models"
	"gitlab.com/PeeK1e/file-drop/pkg/util"
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

func DownloadEncryptedFile(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.RequestURI, "/enc/", "", -1)
	log.Printf("Client %s requested FileID %s", r.RemoteAddr, id)

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("ERR: Could not upgrade HTTP connection to WS %s", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("WARN: WS clean close failed %s", err)
		}
	}()

	sha, isEnc, err := models.CreateChallenge(id)
	if !isEnc || err != nil {
		log.Printf("ERR: file doesn't look encrypted, aborting...")
		return
	}

	//create challenge
	challenge := util.GetRandomString(12)
	sha_512 := sha512.New()
	sha_512.Write([]byte(sha + challenge))
	challengeSha := sha_512.Sum(nil)

	_, err = conn.Write([]byte(challenge))
	if err != nil {
		log.Printf("ERR: WS write stream broken")
		return
	}

	challengeResp := []byte("")
	conn.Read(challengeResp)
	_, err = conn.Write([]byte(challenge))
	if err != nil {
		log.Printf("ERR: WS read stream broken")
		return
	}

	if string(challengeResp) != string(challengeSha) {
		log.Printf("WARN: Client response not matching with server challenge")
		return
	}

	conn.Write()

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

	_, err = file.Stat()
	if err != nil {
		log.Printf("File Corrupted %s", err)
		http.Error(w, "File Corrupted", http.StatusInternalServerError)
		return
	}

	//conn.Write(file.)

	//http.ServeContent(w, r, path, fileStat.ModTime(), file)
}

func CreateChallengeKey() (string, string) {
	return "", ""
}
