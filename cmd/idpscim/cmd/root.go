package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/internal/core"
	"github.com/slashdevops/idp-scim-sync/internal/idp"
	"github.com/slashdevops/idp-scim-sync/internal/repository"
	"github.com/slashdevops/idp-scim-sync/internal/scim"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg               config.Config
	logHandler        slog.Handler
	logHandlerOptions *slog.HandlerOptions
	logger            *slog.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "idpscim",
	Version: version.Version,
	Short:   "Sync your AWS Single Sign-On (SSO) with Google Workspace",
	Long: `
Sync your Google Workspace Groups and Users to AWS Single Sign-On using
AWS SSO SCIM API (https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return sync()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if cfg.IsLambda {
		lambda.Start(rootCmd.Execute)
	}
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cfg = config.New()
	cfg.IsLambda = len(os.Getenv("LAMBDA_TASK_ROOT")) > 0

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfg.ConfigFile, "config-file", "c", config.DefaultConfigFile, "configuration file")

	rootCmd.PersistentFlags().BoolVarP(&cfg.Debug, "debug", "d", config.DefaultDebug, "fast way to set the log-level to debug")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogFormat, "log-format", "f", config.DefaultLogFormat, "set the log format")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogLevel, "log-level", "l", config.DefaultLogLevel, "set the log level [panic|fatal|error|warn|info|debug|trace]")

	rootCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMAccessToken, "aws-scim-access-token", "t", "", "AWS SSO SCIM API Access Token")
	rootCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMAccessTokenSecretName,
		"aws-scim-access-token-secret-name", "j", config.DefaultAWSSCIMAccessTokenSecretName,
		"AWS Secrets Manager secret name for AWS SSO SCIM API Access Token",
	)

	rootCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMEndpoint, "aws-scim-endpoint", "e", "", "AWS SSO SCIM API Endpoint")
	rootCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMEndpointSecretName,
		"aws-scim-endpoint-secret-name", "n", config.DefaultAWSSCIMEndpointSecretName,
		"AWS Secrets Manager secret name for AWS SSO SCIM API Endpoint",
	)

	rootCmd.PersistentFlags().StringVarP(&cfg.AWSS3BucketName, "aws-s3-bucket-name", "b", "", "AWS S3 Bucket name to store the state")
	rootCmd.PersistentFlags().StringVarP(&cfg.AWSS3BucketKey, "aws-s3-bucket-key", "k", config.DefaultAWSS3BucketKey, "AWS S3 Bucket key to store the state")

	rootCmd.PersistentFlags().StringVarP(&cfg.GWSServiceAccountFile,
		"gws-service-account-file", "s", config.DefaultGWSServiceAccountFile,
		"Google Workspace service account file",
	)
	rootCmd.PersistentFlags().StringVarP(&cfg.GWSServiceAccountFileSecretName,
		"gws-service-account-file-secret-name", "o", config.DefaultGWSServiceAccountFileSecretName,
		"AWS Secrets Manager secret name for Google Workspace service account file",
	)

	rootCmd.PersistentFlags().StringVarP(&cfg.GWSUserEmail,
		"gws-user-email", "u", "",
		"GWS user email with allowed access to the Google Workspace Service Account",
	)
	rootCmd.PersistentFlags().StringVarP(&cfg.GWSUserEmailSecretName,
		"gws-user-email-secret-name", "p", config.DefaultGWSUserEmailSecretName,
		"AWS Secrets Manager secret name for GWS user email with allowed access to the Google Workspace Service Account",
	)

	rootCmd.Flags().StringSliceVarP(
		&cfg.GWSGroupsFilter, "gws-groups-filter", "q", []string{""},
		"GWS Groups query parameter, example: --gws-groups-filter 'name:Admin* email:admin*' --gws-groups-filter 'name:Power* email:power*'",
	)

	rootCmd.PersistentFlags().StringVarP(&cfg.SyncMethod, "sync-method", "m", config.DefaultSyncMethod, "Sync method to use [groups]")
	rootCmd.PersistentFlags().BoolVarP(&cfg.UseSecretsManager, "use-secrets-manager", "g", config.DefaultUseSecretsManager, "use AWS Secrets Manager content or not (default false)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
			slog.Error("cannot bind environment variable", "error", err)
			os.Exit(1)
		}
	}

	// when use a lambda, we need to read the config from the environment only
	// so, this is to read the config from file
	if !cfg.IsLambda {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)

		currentDir, err := os.Getwd()
		cobra.CheckErr(err)
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

	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("cannot unmarshal config", "error", err)
	}

	switch strings.ToLower(cfg.LogFormat) {
	case "json":
		logHandler = slog.NewJSONHandler(os.Stdout, logHandlerOptions)
	case "text":
		logHandler = slog.NewTextHandler(os.Stdout, logHandlerOptions)
	default:
		slog.Warn("unknown log format, using text", "format", cfg.LogFormat)
		logHandler = slog.NewTextHandler(os.Stdout, logHandlerOptions)
	}

	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}
	case "info":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelInfo}
	case "warn":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelWarn}
	case "error":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelError, AddSource: true}
	default:
		slog.Warn("unknown log level, setting it to info", "level", cfg.LogLevel)
	}

	if cfg.Debug {
		cfg.LogLevel = "debug"
	}

	if cfg.IsLambda || cfg.UseSecretsManager {
		getSecrets()
	}

	// not implemented yet block
	if cfg.SyncMethod != "groups" {
		slog.Error("only 'sync-method=groups' are implemented")
		os.Exit(1)
	}
}

func getSecrets() {
	slog.Info("reading secrets from AWS Secrets Manager")

	awsConf, err := aws.NewDefaultConf(context.Background())
	if err != nil {
		slog.Error("cannot load aws config", "error", err)
		os.Exit(1)
	}

	svc := secretsmanager.NewFromConfig(awsConf)

	secrets, err := aws.NewSecretsManagerService(svc)
	if err != nil {
		slog.Error("cannot create aws secrets manager service", "error", err)
		os.Exit(1)
	}

	slog.Debug("reading secret", "name", cfg.GWSUserEmailSecretName)
	unwrap, err := secrets.GetSecretValue(context.Background(), cfg.GWSUserEmailSecretName)
	if err != nil {
		slog.Error("cannot get secretmanager value", "error", err)
		os.Exit(1)
	}
	cfg.GWSUserEmail = unwrap

	slog.Debug("reading secret", "name", cfg.GWSServiceAccountFileSecretName)
	unwrap, err = secrets.GetSecretValue(context.Background(), cfg.GWSServiceAccountFileSecretName)
	if err != nil {
		slog.Error("cannot get secretmanager value", "error", err)
		os.Exit(1)
	}
	cfg.GWSServiceAccountFile = unwrap

	slog.Debug("reading secret", "name", cfg.AWSSCIMAccessTokenSecretName)
	unwrap, err = secrets.GetSecretValue(context.Background(), cfg.AWSSCIMAccessTokenSecretName)
	if err != nil {
		slog.Error("cannot get secretmanager value", "error", err)
		os.Exit(1)
	}
	cfg.AWSSCIMAccessToken = unwrap

	slog.Debug("reading secret", "name", cfg.AWSSCIMEndpointSecretName)
	unwrap, err = secrets.GetSecretValue(context.Background(), cfg.AWSSCIMEndpointSecretName)
	if err != nil {
		slog.Error("cannot get secretmanager value", "error", err)
		os.Exit(1)
	}
	cfg.AWSSCIMEndpoint = unwrap
}

func sync() error {
	slog.Debug("viper config", "config", viper.AllSettings())

	if cfg.SyncMethod != "groups" {
		slog.Error("only 'sync-method=groups' are implemented")
		return fmt.Errorf("unknown sync method: %s", cfg.SyncMethod)
	}

	return syncGroups()
}

func syncGroups() error {
	slog.Info("starting sync groups", "codeVersion", version.Version)
	timeStart := time.Now()

	// cfg.GWSServiceAccountFile could be a file path or a content of the file
	gwsServiceAccountContent := []byte(cfg.GWSServiceAccountFile)

	if !cfg.IsLambda {
		gwsServiceAccount, err := os.ReadFile(cfg.GWSServiceAccountFile)
		if err != nil {
			slog.Error("cannot read service account file", "error", err)
		}
		gwsServiceAccountContent = gwsServiceAccount
	}

	gwsAPIScopes := []string{
		"https://www.googleapis.com/auth/admin.directory.group.readonly",
		"https://www.googleapis.com/auth/admin.directory.group.member.readonly",
		"https://www.googleapis.com/auth/admin.directory.user.readonly",
	}

	ctx := context.Background()

	// Google Client Service
	gwsService, err := google.NewService(ctx, cfg.GWSUserEmail, gwsServiceAccountContent, gwsAPIScopes...)
	if err != nil {
		return errors.Wrap(err, "cannot create google service")
	}

	// Google Directory Service
	gwsDS, err := google.NewDirectoryService(gwsService)
	if err != nil {
		return errors.Wrap(err, "cannot create google directory service")
	}

	// Identity Provider Service
	idpService, err := idp.NewIdentityProvider(gwsDS)
	if err != nil {
		return errors.Wrap(err, "cannot create identity provider service")
	}

	// httpClient
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.RetryWaitMin = time.Millisecond * 100

	// set the logger only in debug mode
	if cfg.Debug {
		retryClient.Logger = logger
	} else {
		retryClient.Logger = nil
	}

	httpClient := retryClient.StandardClient()

	// AWS SCIM Service
	awsSCIM, err := aws.NewSCIMService(httpClient, cfg.AWSSCIMEndpoint, cfg.AWSSCIMAccessToken)
	if err != nil {
		return errors.Wrap(err, "cannot create aws scim service")
	}
	awsSCIM.UserAgent = "idp-scim-sync/" + version.Version

	scimService, err := scim.NewProvider(awsSCIM)
	if err != nil {
		return errors.Wrap(err, "cannot create scim provider")
	}

	awsConf, err := aws.NewDefaultConf(context.Background())
	if err != nil {
		slog.Error("cannot load aws config", "error", err)
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(awsConf)
	repo, err := repository.NewS3Repository(s3Client, repository.WithBucket(cfg.AWSS3BucketName), repository.WithKey(cfg.AWSS3BucketKey))
	if err != nil {
		slog.Error("cannot create s3 repository", "error", err)
		os.Exit(1)
	}

	ss, err := core.NewSyncService(idpService, scimService, repo, core.WithIdentityProviderGroupsFilter(cfg.GWSGroupsFilter))
	if err != nil {
		return errors.Wrap(err, "cannot create sync service")
	}

	slog.Debug("app config", "config", cfg)

	if err := ss.SyncGroupsAndTheirMembers(ctx); err != nil {
		return errors.Wrap(err, "cannot sync groups and their members")
	}

	slog.Info("sync groups completed", "duration", time.Since(timeStart).String())

	return nil
}
