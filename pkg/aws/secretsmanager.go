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

var (
	ErrSMClientNil               = errors.New("AWS SecretsManager Client is nil")
	ErrGettingSecretValue        = errors.New("Error getting secret value")
	ErrDecodingBinarySecretValue = errors.New("Error decoding binary secret value")
)

// https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/aws/secretsmanager_mocks.go -source=secretsmanager.go SecretsManagerClientAPI

type SecretsManagerClientAPI interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type SecretsManager struct {
	svc SecretsManagerClientAPI
}

func NewSecretsManagerService(svc SecretsManagerClientAPI) (*SecretsManager, error) {
	if svc == nil {
		return nil, errors.Wrapf(ErrSMClientNil, "NewSecretsManagerService")
	}

	return &SecretsManager{
		svc: svc,
	}, nil
}

func (s *SecretsManager) GetSecretValue(ctx context.Context, secretKey string) (string, error) {
	vIn := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretKey),
		VersionStage: aws.String("AWSCURRENT"),
	}

	r, err := s.svc.GetSecretValue(ctx, vIn)
	if err != nil {
		return "", errors.Wrapf(ErrGettingSecretValue, "GetSecretValue")
	}

	fmt.Printf("len= %v", len(r.SecretBinary))

	var secretString string

	if r.SecretString != nil {
		secretString = *r.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(r.SecretBinary)))
		l, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, r.SecretBinary)
		if err != nil {
			return "", errors.Wrapf(ErrDecodingBinarySecretValue, "GetSecretValue")
		}
		secretString = string(decodedBinarySecretBytes[:l])
	}

	return secretString, nil
}
