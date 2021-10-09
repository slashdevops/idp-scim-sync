package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

// Implement model.Repository interface
// Consume s3.Client

var (
	ErrS3ClientNil         = errors.New("AWS S3 Client is nil")
	ErrGettingS3Object     = errors.New("Error getting S3 object")
	ErrDecodingS3Object    = errors.New("Error decoding S3 object")
	ErrMarshallingState    = errors.New("Error marshalling state")
	ErrPuttingS3Object     = errors.New("Error putting S3 object")
	ErrOptionWithBucketNil = errors.New("Option WithBucket is nil")
	ErrOptionWithKeyNil    = errors.New("Option WithKey is nil")
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../../mocks/repository/s3_mocks.go -source=s3.go S3ClientAPI

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
		return nil, errors.Wrapf(ErrS3ClientNil, "NewS3Repository")
	}

	s3r := &S3Repository{
		mu:     &sync.RWMutex{},
		client: client,
	}

	for _, opt := range opts {
		opt(s3r)
	}

	if s3r.bucket == "" {
		return nil, errors.Wrapf(ErrOptionWithBucketNil, "NewS3Repository")
	}

	if s3r.key == "" {
		return nil, errors.Wrapf(ErrOptionWithKeyNil, "NewS3Repository")
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
		return nil, errors.Wrapf(ErrGettingS3Object, "GetState")
	}
	defer resp.Body.Close()

	log.Debugf("GetObjectOutput: %v", resp)

	var state model.State
	dec := json.NewDecoder(resp.Body)

	if err = dec.Decode(&state); err != nil {
		return nil, errors.Wrapf(ErrDecodingS3Object, "GetState")
	}

	return &state, nil
}

func (r *S3Repository) UpdateState(ctx context.Context, state *model.State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	jsonPayload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return errors.Wrapf(ErrMarshallingState, "UpdateState")
	}

	output, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(r.key),
		Body:   bytes.NewReader(jsonPayload),
	})
	if err != nil {
		return errors.Wrapf(ErrPuttingS3Object, "UpdateState")
	}

	log.Debugf("PutObjectOutput: %v", output)

	return nil
}
