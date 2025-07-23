package aws

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func NewDefaultConf(ctx context.Context) (cfg aws.Config, err error) {
	var confOptions []func(*config.LoadOptions) error

	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		slog.Debug("Using AWS Profile", "profile", profile)
		confOptions = append(confOptions,
			config.WithSharedConfigProfile(profile),
			config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
				options.TokenProvider = stscreds.StdinTokenProvider
			}),
		)
	}

	awsConf, err := config.LoadDefaultConfig(ctx, confOptions...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return awsConf, nil
}
