package routes

import (
	"net/http"

	"gitlab.com/PeeK1e/file-drop/pkg/db"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	if nil != db.GetInstance().Ping() {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
