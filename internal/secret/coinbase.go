package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const (
	versionStage  = "AWSCURRENT"
	envRegion     = "AWS_REGION"
	envSecretName = "COINBASE_SECRET_NAME"
)

type CoinbaseSecretResponse struct {
	Name       string `json:"name"`
	PrivateKey string `json:"privateKey"`
}

func GetCoinbaseSecret(ctx context.Context) (*CoinbaseSecretResponse, error) {
	region, ok := os.LookupEnv(envRegion)
	if !ok {
		return nil, fmt.Errorf("secret: missing %s environment variable", envRegion)
	}

	secretName, ok := os.LookupEnv(envSecretName)
	if !ok {
		return nil, fmt.Errorf("secret: missing %s environment variable", envSecretName)
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)

	resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String(versionStage),
	})
	if err != nil {
		return nil, err
	}

	var secret CoinbaseSecretResponse
	if err := json.Unmarshal([]byte(*resp.SecretString), &secret); err != nil {
		return nil, err
	}

	return &secret, nil
}
