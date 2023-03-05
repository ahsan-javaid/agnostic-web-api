package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"os"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	Expires  int64  `json:"exp"`
}

func GenerateToken(id string) (string, error) {
	sk := os.Getenv("JWT_SECRET")
	jwtSecret := []byte(sk)

	claims := Claims{
		Username: id,
		Expires:  time.Now().Add(time.Hour * 24).Unix(),
  }

	header := map[string]string{
			"alg": "HS256",
			"typ": "JWT",
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
			return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
			return "", err
	}

	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsBase64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signatureInput := fmt.Sprintf("%s.%s", headerBase64, claimsBase64)

	mac := hmac.New(sha256.New, jwtSecret)
	mac.Write([]byte(signatureInput))
	signature := mac.Sum(nil)
	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	token := fmt.Sprintf("%s.%s.%s", headerBase64, claimsBase64, signatureBase64)
	return token, nil
}


func parseToken(tokenString string) (*Claims, error) {
	sk := os.Getenv("JWT_SECRET")
	jwtSecret := []byte(sk)

	tokenParts := strings.Split(tokenString, ".")
	if len(tokenParts) != 3 {
		return nil, fmt.Errorf("Invalid token format")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(tokenParts[0])
	if err != nil {
		return nil, err
	}
	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, err
	}
	if header["alg"] != "HS256" {
		return nil, fmt.Errorf("Unsupported signing algorithm")
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return nil, err
	}
	var tokenClaims Claims
	if err := json.Unmarshal(claimsBytes, &tokenClaims); err != nil {
		return nil, err
	}

	signatureBytes, err := base64.RawURLEncoding.DecodeString(tokenParts[2])
	if err != nil {
		return nil, err
	}
	hmac256 := hmac.New(sha256.New, jwtSecret)
	hmac256.Write([]byte(tokenParts[0] + "." + tokenParts[1]))
	expectedSignature := hmac256.Sum(nil)
	if !hmac.Equal(signatureBytes, expectedSignature) {
		return nil, fmt.Errorf("Invalid signature")
	}

	return &tokenClaims, nil
}

