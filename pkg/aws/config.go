package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	log "github.com/sirupsen/logrus"
)

func NewDefaultConf(ctx context.Context) (cfg aws.Config, err error) {
	var confOptions []func(*config.LoadOptions) error

	if len(os.Getenv("AWS_PROFILE")) > 0 {
		log.Debugf("Using AWS Profile: %s", os.Getenv("AWS_PROFILE"))
		confOptions = append(confOptions,
			config.WithSharedConfigProfile(os.Getenv("AWS_PROFILE")),
			config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
				options.TokenProvider = stscreds.StdinTokenProvider
			}),
		)
	}

	awsConf, err := config.LoadDefaultConfig(
		ctx,
		confOptions...,
	)

	return awsConf, err
}
