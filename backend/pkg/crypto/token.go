package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

const AccessTokenByteLength = 32

func cryptoRandomBytes(size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	return buf, err
}

func GenerateAccessToken() (string, error) {
	b, err := cryptoRandomBytes(AccessTokenByteLength)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
