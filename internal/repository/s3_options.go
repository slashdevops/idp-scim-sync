package repository

// S3RepositoryOption is a function that can be used to configure a S3Repository
// using the functional options pattern.
type S3RepositoryOption func(*S3Repository)

// WithBucket sets the name for the S3 Bucket.
func WithBucket(bucket string) S3RepositoryOption {
	return func(r *S3Repository) {
		r.bucket = bucket
	}
}

// WithKey sets the key for the S3 Bucket.
func WithKey(key string) S3RepositoryOption {
	return func(r *S3Repository) {
		r.key = key
	}
}
