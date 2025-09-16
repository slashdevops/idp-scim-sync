package aws

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"
)

// consume secretsmanager.Client

// ErrSecretManagerClientNil is returned when the SecretsManagerClientAPI is nil.
var ErrSecretManagerClientNil = errors.New("aws: AWS SecretsManager Client cannot be nil")

// https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

//go:generate go tool mockgen -package=mocks -destination=../../mocks/aws/secretsmanager_mocks.go -source=secretsmanager.go SecretsManagerClientAPI

// SecretsManagerClientAPI is the interface to consume the secretsmanager client methods.
type SecretsManagerClientAPI interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// SecretsManagerService is the wrapper for the AWS SecretsManager client.
type SecretsManagerService struct {
	svc SecretsManagerClientAPI
}

// NewSecretsManagerService returns a new SecretsManagerService.
func NewSecretsManagerService(svc SecretsManagerClientAPI) (*SecretsManagerService, error) {
	if svc == nil {
		return nil, ErrSecretManagerClientNil
	}

	return &SecretsManagerService{
		svc: svc,
	}, nil
}

// GetSecretValue returns the secret value for the given secret name or arn.
func (s *SecretsManagerService) GetSecretValue(ctx context.Context, secretKey string) (string, error) {
	vIn := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretKey),
		VersionStage: aws.String("AWSCURRENT"),
	}

	r, err := s.svc.GetSecretValue(ctx, vIn)
	if err != nil {
		return "", fmt.Errorf("aws: error getting secret value: %v", err)
	}

	var secretString string

	if r.SecretString != nil {
		secretString = *r.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(r.SecretBinary)))
		l, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, r.SecretBinary)
		if err != nil {
			return "", fmt.Errorf("aws: error decoding secret binary value: %v", err)
		}
		secretString = string(decodedBinarySecretBytes[:l])
	}

	return secretString, nil
}
