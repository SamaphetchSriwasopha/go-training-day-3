package shorten

import (
	"crypto/rand"
	"encoding/base64"
)

func NewShorter() shortenFunc {
	return shortenFunc(shortenURL)
}

func shortenURL() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil // e.g. "aZ3xQ1"
}

type shortenFunc func() (string, error)

func (fn shortenFunc) generate() (string, error) {
	return fn()
}
