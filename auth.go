package fonbnk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (fc *FonbnkClient) setAuthHeaders(req *http.Request) error {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	signature, err := generateSignature(fc.clientSecret, timestamp, req.URL.RequestURI())
	log.Println(req.URL.RequestURI())
	if err != nil {
		return err
	}

	req.Header.Set("x-client-id", fc.clientID)
	req.Header.Set("x-timestamp", timestamp)
	req.Header.Set("x-signature", signature)

	return nil
}

func generateSignature(clientSecret string, timestamp string, endpoint string) (string, error) {
	decodedSecret, err := base64.RawStdEncoding.DecodeString(clientSecret)
	if err != nil {
		return "", fmt.Errorf("failed to decode client secret: %v", err)
	}

	h := hmac.New(sha256.New, decodedSecret)
	h.Write([]byte(timestamp + ":" + endpoint))

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
