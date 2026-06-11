package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

func SHA256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func HMACSHA256(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func RandomToken(bytesLen int) (string, error) {
	if bytesLen <= 0 {
		bytesLen = 32
	}
	b := make([]byte, bytesLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Password(password string) (string, error) {
	salt, err := RandomToken(16)
	if err != nil {
		return "", err
	}
	sum := SHA256(salt + ":" + password)
	return fmt.Sprintf("sha256$%s$%s", salt, sum), nil
}

func CheckPassword(encoded, password string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 3 || parts[0] != "sha256" {
		return false
	}
	expected := SHA256(parts[1] + ":" + password)
	return hmac.Equal([]byte(parts[2]), []byte(expected))
}

func RequireToken(n int) string {
	t, err := RandomToken(n)
	if err != nil {
		panic(errors.Join(errors.New("hash: random token failed"), err))
	}
	return t
}
