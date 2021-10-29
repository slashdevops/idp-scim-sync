package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

// Implement model.Repository interface
// Consume s3.Client

var (
	ErrS3ClientNil         = errors.New("s3: AWS S3 Client is nil")
	ErrOptionWithBucketNil = errors.New("s3: option WithBucket is nil")
	ErrOptionWithKeyNil    = errors.New("s3: option WithKey is nil")
	ErrStateNil            = errors.New("s3: state is nil")
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../../mocks/repository/s3_mocks.go -source=s3.go S3ClientAPI

type S3ClientAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3Repository struct {
	mu     *sync.RWMutex
	bucket string
	key    string
	client S3ClientAPI
}

func NewS3Repository(client S3ClientAPI, opts ...S3RepositoryOption) (*S3Repository, error) {
	if client == nil {
		return nil, ErrS3ClientNil
	}

	s3r := &S3Repository{
		mu:     &sync.RWMutex{},
		client: client,
	}

	for _, opt := range opts {
		opt(s3r)
	}

	if s3r.bucket == "" {
		return nil, ErrOptionWithBucketNil
	}

	if s3r.key == "" {
		return nil, ErrOptionWithKeyNil
	}

	return s3r, nil
}

func (r *S3Repository) GetState(ctx context.Context) (*model.State, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resp, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(r.key),
	})
	if err != nil {
		return nil, fmt.Errorf("s3: error getting S3 object: %w", err)
	}
	defer resp.Body.Close()

	log.Debugf("GetObjectOutput: %s", utils.ToJSON(resp))

	var state model.State
	dec := json.NewDecoder(resp.Body)

	if err = dec.Decode(&state); err != nil {
		return nil, fmt.Errorf("s3: error decoding S3 object: %w", err)
	}

	return &state, nil
}

func (r *S3Repository) SetState(ctx context.Context, state *model.State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if state == nil {
		return ErrStateNil
	}

	jsonPayload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("s3: error marshalling state: %w", err)
	}

	output, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(r.key),
		Body:   bytes.NewReader(jsonPayload),
	})
	if err != nil {
		return fmt.Errorf("s3: error putting S3 object: %w", err)
	}

	log.Debugf("PutObjectOutput: %s", utils.ToJSON(output))

	return nil
}
