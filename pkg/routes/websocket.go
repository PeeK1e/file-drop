package routes

import (
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type socket struct {
	C net.Conn
}

func (s *socket) NewSocket(w http.ResponseWriter, r *http.Request) bool {
	c, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("ERR: WS Upgrade Failed")
		return false
	}
	s.C = c
	return true
}

func (s *socket) WriteMessage(m string) bool {
	if err := wsutil.WriteClientMessage(s.C, ws.OpText, []byte(m)); err != nil {
		log.Printf("WARN: Writing message failed")
		return false
	}
	return true
}

func (s *socket) WriteBytes(b []byte) bool {
	if err := wsutil.WriteClientMessage(s.C, ws.OpBinary, b); err != nil {
		log.Printf("WARN: Writing bytes failed")
		return false
	}
	return true
}

// reads from websocket and returns data as []byte ok if operation succeded
func (s *socket) ReadClientData(m *string) bool {
	data, opcode, err := wsutil.ReadClientData(s.C)
	if err != nil {
		log.Printf("WARN: Read message failed: %s", err)
		return false
	}

	if opcode != ws.OpBinary && opcode != ws.OpText {
		log.Printf("WARN: Data not in Text or Binary format")
		return false
	}

	*m = string(data)
	return true
}

func (s *socket) Close() {
	if err := s.C.Close(); err != nil {
		log.Printf("WARN: Could not close ws cleanly: %s", err)
	}
}
