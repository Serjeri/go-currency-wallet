package lib

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashedPassword(password string) string {
	hashedPassword := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashedPassword[:])
}
