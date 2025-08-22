package coinbase

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (c *Client) BuildJWT(method, path string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": c.keyName,
		"iss": "cdp",
		"nbf": jwt.NewNumericDate(now),
		"exp": jwt.NewNumericDate(now.Add(2 * time.Minute)),
		"uri": fmt.Sprintf("%s %s%s", method, c.baseURL, path),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = c.keyName
	token.Header["nonce"] = uuid.New().String()

	signedToken, err := token.SignedString(c.signer)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
