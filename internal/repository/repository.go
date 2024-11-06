package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -package=mocks -destination=../../mocks/repository/repository_mocks.go -source=repository.go

// S3ClientAPI is an interface to consume S3 client methods
type S3ClientAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}
