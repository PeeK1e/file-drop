package util

import "math/rand"

const (
	letters = string("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func GetRandomString(length int) string {
	str := ""
	for i := 0; i < length; i++ {
		str += string(letters[rand.Intn(len(letters))])
	}
	return str
}
