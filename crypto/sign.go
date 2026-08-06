package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	HeaderAppID     = "X-App-Id"
	HeaderTimestamp = "X-Timestamp"
	HeaderNonce     = "X-Nonce"
	HeaderSignature = "X-Signature"

	ExpireSeconds = 300
)

type SignHeader struct {
	AppID     string
	Timestamp int64
	Nonce     string
	Signature string
}

func SHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func BuildSignString(
	method string,
	path string,
	timestamp int64,
	nonce string,
	body []byte,
) string {

	bodyHash := SHA256Hex(body)

	return fmt.Sprintf(
		"%s\n%s\n%d\n%s\n%s",
		method,
		path,
		timestamp,
		nonce,
		bodyHash,
	)
}

func Sign(privatePEM string, signString string) (string, error) {

	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return "", errors.New("private key error")
	}

	keyAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	key := keyAny.(*rsa.PrivateKey)

	hash := sha256.Sum256([]byte(signString))

	signature, err := rsa.SignPKCS1v15(
		rand.Reader,
		key,
		crypto.SHA256,
		hash[:],
	)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func Verify(publicPEM string, signString string, sign string) error {

	block, _ := pem.Decode([]byte(publicPEM))
	if block == nil {
		return errors.New("public key error")
	}

	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	pub := pubAny.(*rsa.PublicKey)

	hash := sha256.Sum256([]byte(signString))

	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(
		pub,
		crypto.SHA256,
		hash[:],
		signature,
	)
}

func VerifyTimestamp(timestamp int64) error {

	now := time.Now().Unix()

	if timestamp < now-ExpireSeconds {
		return errors.New("request expired")
	}

	if timestamp > now+ExpireSeconds {
		return errors.New("invalid timestamp")
	}

	return nil
}

func ParseHeader(r *http.Request) (*SignHeader, error) {

	ts, err := strconv.ParseInt(
		r.Header.Get(HeaderTimestamp),
		10,
		64,
	)

	if err != nil {
		return nil, err
	}

	return &SignHeader{
		AppID:     r.Header.Get(HeaderAppID),
		Timestamp: ts,
		Nonce:     r.Header.Get(HeaderNonce),
		Signature: r.Header.Get(HeaderSignature),
	}, nil
}
