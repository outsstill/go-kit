package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	HeaderAccessKey = "X-Tunnel-Access-Key"
	HeaderTimestamp = "X-Tunnel-Timestamp"
	HeaderNonce     = "X-Tunnel-Nonce"
	HeaderSignature = "X-Tunnel-Signature"

	// 请求有效时间（秒）
	ExpireSeconds = 300
)

type Header struct {
	AccessKey string
	Timestamp int64
	Nonce     string
	Signature string
}

func SHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func HmacBuildSignString(
	method string,
	path string,
	timestamp int64,
	nonce string,
	body []byte,
) string {

	return fmt.Sprintf(
		"%s\n%s\n%d\n%s\n%s",
		strings.ToUpper(method),
		path,
		timestamp,
		nonce,
		SHA256(body),
	)
}

func HmacSign(secret string, signString string) string {

	h := hmac.New(
		sha256.New,
		[]byte(secret),
	)

	h.Write([]byte(signString))

	return hex.EncodeToString(h.Sum(nil))
}

func HmacVerify(secret string, signString string, signature string) bool {

	serverSign := HmacSign(secret, signString)

	return hmac.Equal(
		[]byte(serverSign),
		[]byte(signature),
	)
}

func HmacVerifyTimestamp(timestamp int64) error {

	now := time.Now().Unix()

	if timestamp < now-ExpireSeconds {
		return errors.New("request expired")
	}

	if timestamp > now+ExpireSeconds {
		return errors.New("invalid timestamp")
	}

	return nil
}

func HmacParseHeader(
	accessKey string,
	timestamp string,
	nonce string,
	signature string,
) (*Header, error) {

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, err
	}

	return &Header{
		AccessKey: accessKey,
		Timestamp: ts,
		Nonce:     nonce,
		Signature: signature,
	}, nil
}
