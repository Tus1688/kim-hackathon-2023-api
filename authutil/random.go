package authutil

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString generates base64 url encoded random string with certain length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
