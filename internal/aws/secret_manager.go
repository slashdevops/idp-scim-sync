package aws

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretsManagerService interface {
	GetSecretValue(secretKey string) (string, error)
}

type secretsManager struct {
	svc *secretsmanager.SecretsManager
}

func NewSecretsManagerService(svc *secretsmanager.SecretsManager) SecretsManagerService {
	return &secretsManager{
		svc: svc,
	}
}

func (s *secretsManager) GetSecretValue(secretKey string) (string, error) {

	r, err := s.svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
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
