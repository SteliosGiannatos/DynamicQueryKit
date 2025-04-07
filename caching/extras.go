package caching

import (
	"crypto/sha1"
	"encoding/hex"
)

func hashKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	hashBytes := h.Sum(nil)
	key = hex.EncodeToString(hashBytes)
	return key
}
