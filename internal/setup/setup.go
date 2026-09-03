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
	"github.com/slashdevops/httpx"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/internal/core"
	"github.com/slashdevops/idp-scim-sync/internal/idp"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/repository"
	"github.com/slashdevops/idp-scim-sync/internal/scim"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
)

// Logger builds a slog.Logger for the given level and format, installs it as
// the process default, and returns it.
//
// Both binaries route through this function so their logging behaviour cannot
// drift apart.
func Logger(logLevel, logFormat string) *slog.Logger {
	logHandlerOptions := handlerOptions(logLevel)

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

// handlerOptions maps a configured log level onto slog handler options.
//
// "fatal" and "panic" are accepted because config.Validate has always allowed
// them — they predate the move from logrus to log/slog, which has no equivalent
// levels. Rather than silently degrading them to info (which is what the old
// default branch did, complete with a misleading "unknown log level" warning),
// they map to slog.LevelError, the closest slog equivalent. Existing
// deployments configured with either value keep working and now get the
// severity they asked for.
func handlerOptions(logLevel string) *slog.HandlerOptions {
	switch strings.ToLower(logLevel) {
	case "debug":
		return &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}
	case "info":
		return &slog.HandlerOptions{Level: slog.LevelInfo}
	case "warn":
		return &slog.HandlerOptions{Level: slog.LevelWarn}
	case "error", "fatal", "panic":
		return &slog.HandlerOptions{Level: slog.LevelError, AddSource: true}
	default:
		slog.Warn("unknown log level, setting it to info", "level", logLevel)
		return &slog.HandlerOptions{Level: slog.LevelInfo}
	}
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
		"sync_user_fields",
	}
	for _, e := range envVars {
		if err := viper.BindEnv(e); err != nil {
			return fmt.Errorf("cannot bind environment variable: %w", err)
		}
	}

	// when use a lambda, we need to read the config from the environment only
	// so, this is to read the config from file
	if !cfg.IsLambda {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot get user home directory: %w", err)
		}
		viper.AddConfigPath(home)

		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get current directory: %w", err)
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
		return fmt.Errorf("cannot unmarshal config: %w", err)
	}

	if cfg.Debug {
		cfg.LogLevel = "debug"
	}

	return nil
}

// secretsFetcher reads a secret value by name. It is the subset of
// aws.SecretsManagerService that Secrets depends on, declared here so the
// concurrent fetch below can be exercised without an AWS client.
type secretsFetcher interface {
	GetSecretValue(ctx context.Context, secretKey string) (string, error)
}

// Secrets populates the secret-backed fields of cfg from AWS Secrets Manager.
//
// ctx governs the whole operation: it is used to load the AWS configuration and
// is passed through to every secret read, so a Lambda nearing its deadline can
// actually abandon the requests instead of leaving them running until the
// runtime freezes the environment.
func Secrets(ctx context.Context, cfg *config.Config) error {
	slog.Info("reading secrets from AWS Secrets Manager")

	awsConf, err := aws.NewDefaultConf(ctx)
	if err != nil {
		return fmt.Errorf("cannot load aws config: %w", err)
	}

	svc := secretsmanager.NewFromConfig(awsConf)

	secrets, err := aws.NewSecretsManagerService(svc)
	if err != nil {
		return fmt.Errorf("cannot create aws secrets manager service: %w", err)
	}

	return fetchSecrets(ctx, secrets, cfg)
}

// fetchSecrets reads the four secrets concurrently and assigns each to its
// field in cfg.
//
// The four goroutines write to four distinct fields, so no lock is needed. They
// share one errgroup context, so the first failure cancels whichever reads are
// still outstanding rather than letting them run to completion for a result the
// caller will discard.
func fetchSecrets(ctx context.Context, secrets secretsFetcher, cfg *config.Config) error {
	reads := []struct {
		name string
		dest *string
	}{
		{cfg.GWSUserEmailSecretName, &cfg.GWSUserEmail},
		{cfg.GWSServiceAccountFileSecretName, &cfg.GWSServiceAccountFile},
		{cfg.AWSSCIMAccessTokenSecretName, &cfg.AWSSCIMAccessToken},
		{cfg.AWSSCIMEndpointSecretName, &cfg.AWSSCIMEndpoint},
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, read := range reads {
		g.Go(func() error {
			slog.Debug("reading secret", "name", read.name)

			value, err := secrets.GetSecretValue(ctx, read.name)
			if err != nil {
				return fmt.Errorf("cannot get secretmanager value for %q: %w", read.name, err)
			}

			*read.dest = value

			return nil
		})
	}

	return g.Wait()
}

// SyncService sets up the sync service
func SyncService(ctx context.Context, cfg *config.Config) (*core.SyncService, error) {
	// cfg.GWSServiceAccountFile could be a file path or a content of the file
	gwsServiceAccountContent := []byte(cfg.GWSServiceAccountFile)

	if !cfg.IsLambda {
		gwsServiceAccount, err := os.ReadFile(cfg.GWSServiceAccountFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read google workspace service account file: %w", err)
		}
		gwsServiceAccountContent = gwsServiceAccount
	}

	idpClient := httpx.NewClientBuilder().
		WithMaxRetries(10).
		WithRetryStrategy(httpx.ExponentialBackoffStrategy).
		WithRetryBaseDelay(500 * time.Millisecond).
		WithRetryMaxDelay(10 * time.Second).
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
		return nil, fmt.Errorf("cannot create google service: %w", err)
	}

	// Build the sync field set from configuration
	syncFieldSet := model.NewSyncFieldSet(cfg.SyncUserFields)

	// Google Directory Service
	gwsDS, err := google.NewDirectoryService(gwsService, google.WithSyncFieldSet(syncFieldSet))
	if err != nil {
		return nil, fmt.Errorf("cannot create google directory service: %w", err)
	}

	// Identity Provider Service
	idpService, err := idp.NewIdentityProvider(gwsDS, idp.WithSyncFieldSet(syncFieldSet))
	if err != nil {
		return nil, fmt.Errorf("cannot create identity provider service: %w", err)
	}

	// AWS SCIM Service

	// httpClient with jitter backoff to avoid thundering herd on 429 rate limits
	scimClient := httpx.NewClientBuilder().
		WithMaxRetries(10).
		WithRetryStrategy(httpx.JitterBackoffStrategy).
		WithRetryBaseDelay(500 * time.Millisecond).
		WithRetryMaxDelay(10 * time.Second).
		Build()

	awsSCIM, err := aws.NewSCIMService(scimClient, cfg.AWSSCIMEndpoint, cfg.AWSSCIMAccessToken)
	if err != nil {
		return nil, fmt.Errorf("cannot create aws scim service: %w", err)
	}
	awsSCIM.UserAgent = userAgent

	scimService, err := scim.NewProvider(awsSCIM)
	if err != nil {
		return nil, fmt.Errorf("cannot create scim provider: %w", err)
	}

	awsConf, err := aws.NewDefaultConf(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot load aws config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsConf)
	repo, err := repository.NewS3Repository(s3Client, repository.WithBucket(cfg.AWSS3BucketName), repository.WithKey(cfg.AWSS3BucketKey))
	if err != nil {
		return nil, fmt.Errorf("cannot create s3 repository: %w", err)
	}

	ss, err := core.NewSyncService(idpService, scimService, repo, core.WithIdentityProviderGroupsFilter(cfg.GWSGroupsFilter))
	if err != nil {
		return nil, fmt.Errorf("cannot create sync service: %w", err)
	}

	return ss, nil
}
