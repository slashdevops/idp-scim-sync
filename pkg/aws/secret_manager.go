package aws

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/
type SecretsManagerService interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type SecretsManager struct {
	svc SecretsManagerService
}

func NewSecretsManagerService(svc SecretsManagerService) *SecretsManager {
	return &SecretsManager{
		svc: svc,
	}
}

func (s *SecretsManager) GetSecretValue(ctx context.Context, secretKey string) (string, error) {
	r, err := s.svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretKey),
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil {
		return "", err
	}

	var secretString string

	if r.SecretString != nil {
		secretString = *r.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(r.SecretBinary)))
		l, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, r.SecretBinary)
		if err != nil {
			return "", err
		}
		secretString = string(decodedBinarySecretBytes[:l])
	}

	return secretString, nil
}
