package s3

type S3RepositoryOption func(*S3Repository)

func WithBucket(bucket string) S3RepositoryOption {
	return func(r *S3Repository) {
		r.bucket = bucket
	}
}

func WithKey(key string) S3RepositoryOption {
	return func(r *S3Repository) {
		r.key = key
	}
}
