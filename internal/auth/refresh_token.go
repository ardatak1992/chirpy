package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {

	var b []byte
	rand.Read(b)
	tokenString := hex.EncodeToString(b)
	return tokenString
}
