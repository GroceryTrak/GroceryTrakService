package clients

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type CognitoClaims struct {
	Sub string `json:"sub"`
}

func ParseToken(tokenString string) (*CognitoClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	// Decode the claims part (second part of the JWT)
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims CognitoClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
