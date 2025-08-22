package coinbase

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"coinbase-advanced-automation/internal/secret"
)

const baseURL = "api.coinbase.com"

type Client struct {
	baseURL    string
	keyName    string
	signer     crypto.Signer
	httpClient *http.Client
}

func NewClient(coinbaseSecret *secret.CoinbaseSecretResponse) (*Client, error) {
	privateKey := strings.ReplaceAll(coinbaseSecret.PrivateKey, `\n`, "\n")

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("coinbase: failed to decode PEM block for the private key")
	}

	signer, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL: baseURL,
		keyName: coinbaseSecret.Name,
		signer:  signer,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *Client) MakeRequest(ctx context.Context, method, path, signedToken string, body io.Reader) (*http.Response, error) {
	fullURL := fmt.Sprintf("https://%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+signedToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
