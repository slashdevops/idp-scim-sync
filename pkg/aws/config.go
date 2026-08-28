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

// NewDefaultConf loads the default AWS configuration for the current
// environment.
//
// When AWS_PROFILE is set, the named shared-config profile is used and
// stscreds.StdinTokenProvider is wired in so an MFA-protected assume-role
// profile can prompt on stdin. That path is for local CLI use; in Lambda the
// execution role is resolved from the environment and no profile is set.
func NewDefaultConf(ctx context.Context) (aws.Config, error) {
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
