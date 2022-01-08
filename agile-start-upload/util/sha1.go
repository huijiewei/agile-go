package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func sha1Encode(s string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(s))

	return hex.EncodeToString(sha1.Sum(nil))
}

func hmacSha1Encode(s, key string) string {
	hmac := hmac.New(sha1.New, []byte(key))
	hmac.Write([]byte(s))

	return hex.EncodeToString(hmac.Sum(nil))
}
