package util

import "encoding/hex"

func ArrayBufferToString(str []byte) string {
	return hex.EncodeToString(str)
}
