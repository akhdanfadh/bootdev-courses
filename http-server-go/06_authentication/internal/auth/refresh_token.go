package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	data := make([]byte, 32)        // 32 bytes = 256 bits
	rand.Read(data)                 // fills data with random data
	return hex.EncodeToString(data) // convert to hex string
}
