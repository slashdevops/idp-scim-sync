package setup

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/p2p-b2b/httpretrier"
	"github.com/pkg/errors"
	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/internal/core"
	"github.com/slashdevops/idp-scim-sync/internal/idp"
	"github.com/slashdevops/idp-scim-sync/internal/repository"
	"github.com/slashdevops/idp-scim-sync/internal/scim"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
	"github.com/spf13/viper"
)

// Logger sets up the logger
func Logger(logLevel, logFormat string) *slog.Logger {
	var logHandlerOptions *slog.HandlerOptions
	switch strings.ToLower(logLevel) {
	case "debug":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}
	case "info":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelInfo}
	case "warn":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelWarn}
	case "error":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelError, AddSource: true}
	default:
		slog.Warn("unknown log level, setting it to info", "level", logLevel)
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelInfo}
	}

	var logHandler slog.Handler
	switch strings.ToLower(logFormat) {
	case "json":
		logHandler = slog.NewJSONHandler(os.Stdout, logHandlerOptions)
	case "text":
		logHandler = slog.NewTextHandler(os.Stdout, logHandlerOptions)
	default:
		slog.Warn("unknown log format, using text", "format", logFormat)
		logHandler = slog.NewTextHandler(os.Stdout, logHandlerOptions)
	}

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	return logger
}

// Configuration sets up the configuration
func Configuration(cfg *config.Config) error {
	viper.SetEnvPrefix("idpscim") // allow to read in from environment

	envVars := []string{
		"log_level",
		"log_format",
		"sync_method",
		"aws_s3_bucket_name",
		"aws_s3_bucket_key",
		"gws_user_email",
		"gws_user_email_secret_name",
		"gws_service_account_file",
		"gws_service_account_file_secret_name",
		"gws_groups_filter",
		"aws_scim_access_token",
		"aws_scim_access_token_secret_name",
		"aws_scim_endpoint",
		"aws_scim_endpoint_secret_name",
		"use_secrets_manager",
	}
	for _, e := range envVars {
		if err := viper.BindEnv(e); err != nil {
			return errors.Wrap(err, "cannot bind environment variable")
		}
	}

	// when use a lambda, we need to read the config from the environment only
	// so, this is to read the config from file
	if !cfg.IsLambda {
		home, err := os.UserHomeDir()
		if err != nil {
			return errors.Wrap(err, "cannot get user home directory")
		}
		viper.AddConfigPath(home)

		currentDir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "cannot get current directory")
		}
		viper.AddConfigPath(currentDir)

		fileDir := filepath.Dir(cfg.ConfigFile)
		viper.AddConfigPath(fileDir)

		// Search config in home directory with name "downloader" (without extension).
		fileExtension := filepath.Ext(cfg.ConfigFile)
		fileExtensionName := fileExtension[1:]
		viper.SetConfigType(fileExtensionName)

		fileNameExt := filepath.Base(cfg.ConfigFile)
		fileName := fileNameExt[0 : len(fileNameExt)-len(fileExtension)]
		viper.SetConfigName(fileName)

		slog.Debug("configuration file", "dir", fileDir, "name", fileName, "extension", fileExtension)

		if err := viper.ReadInConfig(); err == nil {
			slog.Info("using config file", "file", viper.ConfigFileUsed())
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return errors.Wrap(err, "cannot unmarshal config")
	}

	if cfg.Debug {
		cfg.LogLevel = "debug"
	}

	return nil
}

// Secrets sets up the secrets
func Secrets(cfg *config.Config) error {
	slog.Info("reading secrets from AWS Secrets Manager")

	awsConf, err := aws.NewDefaultConf(context.Background())
	if err != nil {
		return errors.Wrap(err, "cannot load aws config")
	}

	svc := secretsmanager.NewFromConfig(awsConf)

	secrets, err := aws.NewSecretsManagerService(svc)
	if err != nil {
		return errors.Wrap(err, "cannot create aws secrets manager service")
	}

	// create a channel to receive the results
	results := make(chan error, 4)

	go func() {
		slog.Debug("reading secret", "name", cfg.GWSUserEmailSecretName)
		unwrap, err := secrets.GetSecretValue(context.Background(), cfg.GWSUserEmailSecretName)
		if err != nil {
			results <- errors.Wrap(err, "cannot get secretmanager value")
			return
		}
		cfg.GWSUserEmail = unwrap
		results <- nil
	}()

	go func() {
		slog.Debug("reading secret", "name", cfg.GWSServiceAccountFileSecretName)
		unwrap, err := secrets.GetSecretValue(context.Background(), cfg.GWSServiceAccountFileSecretName)
		if err != nil {
			results <- errors.Wrap(err, "cannot get secretmanager value")
			return
		}
		cfg.GWSServiceAccountFile = unwrap
		results <- nil
	}()

	go func() {
		slog.Debug("reading secret", "name", cfg.AWSSCIMAccessTokenSecretName)
		unwrap, err := secrets.GetSecretValue(context.Background(), cfg.AWSSCIMAccessTokenSecretName)
		if err != nil {
			results <- errors.Wrap(err, "cannot get secretmanager value")
			return
		}
		cfg.AWSSCIMAccessToken = unwrap
		results <- nil
	}()

	go func() {
		slog.Debug("reading secret", "name", cfg.AWSSCIMEndpointSecretName)
		unwrap, err := secrets.GetSecretValue(context.Background(), cfg.AWSSCIMEndpointSecretName)
		if err != nil {
			results <- errors.Wrap(err, "cannot get secretmanager value")
			return
		}
		cfg.AWSSCIMEndpoint = unwrap
		results <- nil
	}()

	// wait for all the goroutines to finish
	for i := 0; i < 4; i++ {
		if err := <-results; err != nil {
			return err
		}
	}

	return nil
}

// SyncService sets up the sync service
func SyncService(ctx context.Context, cfg *config.Config) (*core.SyncService, error) {
	// cfg.GWSServiceAccountFile could be a file path or a content of the file
	gwsServiceAccountContent := []byte(cfg.GWSServiceAccountFile)

	if !cfg.IsLambda {
		gwsServiceAccount, err := os.ReadFile(cfg.GWSServiceAccountFile)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read google workspace service account file")
		}
		gwsServiceAccountContent = gwsServiceAccount
	}

	idpClient := httpretrier.NewClientBuilder().
		WithTimeout(30 * time.Second).                        // Overall request timeout
		WithMaxRetries(5).                                    // Retry up to 3 times
		WithRetryStrategy(httpretrier.JitterBackoffStrategy). // Use jitter backoff to avoid thundering herd
		WithRetryBaseDelay(500 * time.Millisecond).           // Start with 500ms delay (httpretrier default)
		WithRetryMaxDelay(5 * time.Second).                   // Cap at 5 seconds
		WithMaxIdleConns(10).                                 // Max idle connections
		WithMaxIdleConnsPerHost(10).                          // Max idle connections per host
		WithIdleConnTimeout(90 * time.Second).                // Idle connection timeout
		Build()

	userAgent := fmt.Sprintf("idp-scim-sync/%s", version.Version)

	gServiceConfig := google.DirectoryServiceConfig{
		UserEmail:      cfg.GWSUserEmail,
		ServiceAccount: gwsServiceAccountContent,
		Scopes:         cfg.GWSServiceAccountScopes,
		UserAgent:      userAgent,
		Client:         idpClient,
	}

	// Google Client Service
	gwsService, err := google.NewService(ctx, gServiceConfig)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create google service")
	}

	// Google Directory Service
	gwsDS, err := google.NewDirectoryService(gwsService)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create google directory service")
	}

	// Identity Provider Service
	idpService, err := idp.NewIdentityProvider(gwsDS)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create identity provider service")
	}

	// AWS SCIM Service

	scimClient := httpretrier.NewClientBuilder().
		WithTimeout(30 * time.Second).                        // Overall request timeout
		WithMaxRetries(10).                                   // Retry up to 3 times
		WithRetryStrategy(httpretrier.JitterBackoffStrategy). // Use jitter backoff to avoid thundering herd
		WithRetryBaseDelay(500 * time.Millisecond).           // Start with 500ms delay (httpretrier default)
		WithRetryMaxDelay(10 * time.Second).                  // Cap at 5 seconds
		WithMaxIdleConns(10).                                 // Max idle connections
		WithMaxIdleConnsPerHost(10).                          // Max idle connections per host
		WithIdleConnTimeout(90 * time.Second).                // Idle connection timeout
		Build()

	awsSCIM, err := aws.NewSCIMService(scimClient, cfg.AWSSCIMEndpoint, cfg.AWSSCIMAccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create aws scim service")
	}
	awsSCIM.UserAgent = userAgent

	scimService, err := scim.NewProvider(awsSCIM)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create scim provider")
	}

	awsConf, err := aws.NewDefaultConf(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "cannot load aws config")
	}

	s3Client := s3.NewFromConfig(awsConf)
	repo, err := repository.NewS3Repository(s3Client, repository.WithBucket(cfg.AWSS3BucketName), repository.WithKey(cfg.AWSS3BucketKey))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create s3 repository")
	}

	ss, err := core.NewSyncService(idpService, scimService, repo, core.WithIdentityProviderGroupsFilter(cfg.GWSGroupsFilter))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create sync service")
	}

	return ss, nil
}
