package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}
