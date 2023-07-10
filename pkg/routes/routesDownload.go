package routes

import (
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
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

	wsChallenge(w, r)

	//conn.Write(file.)

	//http.ServeContent(w, r, path, fileStat.ModTime(), file)
}

func getEncFile(id string) (string, string) {
	sha, isEnc, err := models.GetEncryptionDetails(id)
	if !isEnc || err != nil {
		log.Printf("ERR: file doesn't look encrypted, aborting...")
		return "", ""
	}

	_, path, err := models.GetFileByID(id)
	if err != nil {
		log.Printf("Couldn't retrieve Database Entry %s", err)
		//http.Error(w, "File not found", http.StatusNotFound)
		return "", ""
	}

	return sha, path
}

// func wsChallenge(w http.ResponseWriter, r *http.Request, file *os.File, sha string) {
func wsChallenge(w http.ResponseWriter, r *http.Request) {
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

	fileID, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		log.Printf("ERR: WS read stream broken")
		return
	}
	log.Printf("DEBUG: file id: %s // %s", fileID, string(fileID))

	pwSha, filePath := getEncFile(string(fileID))
	if pwSha == "" {
		log.Printf("ERR: Could not open file")
		return
	}

	//create challenge
	challenge := util.GetRandomString(12)
	sha_512 := sha512.New()
	sha_512.Write([]byte(pwSha + challenge))
	challengeSha := sha_512.Sum(nil)

	err = wsutil.WriteServerMessage(conn, ws.OpText, []byte(challenge))
	if err != nil {
		log.Printf("ERR: WS write stream broken")
		return
	}

	challengeResp, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		log.Printf("ERR: WS read stream broken")
		return
	}

	strChallengeResp := hex.EncodeToString(challengeResp)
	strChallenge := hex.EncodeToString(challengeSha)

	log.Printf("Server challenge: %s, Client Response %s", string(strChallenge), string(strChallengeResp))

	if string(challengeResp) != string(challengeSha) {
		log.Printf("WARN: Client response not matching with server challenge")
		return
	}

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return
	}

	err = wsutil.WriteServerMessage(conn, ws.OpBinary, []byte(fileData))
	if err != nil {
		log.Printf("ERR: WS write stream broken")
		return
	}
}

func CreateChallengeKey() (string, string) {
	return "", ""
}
