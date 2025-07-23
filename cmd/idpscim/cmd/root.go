// Package cmd provides the root command and configuration for the idpscim CLI application.
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/internal/setup"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "idpscim",
	Version: version.Version,
	Short:   "Sync your AWS Single Sign-On (SSO) with Google Workspace",
	Long: `
Sync your Google Workspace Groups and Users to AWS Single Sign-On using
AWS SSO SCIM API (https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
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

	cobra.OnInitialize(func() {
		if err := setup.Configuration(&cfg); err != nil {
			slog.Error("cannot setup configuration", "error", err)
			os.Exit(1)
		}
		setup.Logger(cfg.LogLevel, cfg.LogFormat)

		if cfg.IsLambda || cfg.UseSecretsManager {
			if err := setup.Secrets(&cfg); err != nil {
				slog.Error("cannot get secrets", "error", err)
				os.Exit(1)
			}
		}
	})

	rootCmd.PersistentFlags().StringVarP(&cfg.ConfigFile, "config-file", "c", config.DefaultConfigFile, "configuration file")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Debug, "debug", "d", config.DefaultDebug, "fast way to set the log-level to debug")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogFormat, "log-format", "f", config.DefaultLogFormat, "set the log format")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogLevel, "log-level", "l", config.DefaultLogLevel, "set the log level [panic|fatal|error|warn|info|debug|trace]")
	rootCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMAccessToken, "aws-scim-access-token", "t", "", "AWS SSO SCIM API Access Token")
	rootCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMAccessTokenSecretName, "aws-scim-access-token-secret-name", "j", config.DefaultAWSSCIMAccessTokenSecretName, "AWS Secrets Manager secret name for AWS SSO SCIM API Access Token")
	rootCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMEndpoint, "aws-scim-endpoint", "e", "", "AWS SSO SCIM API Endpoint")
	rootCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMEndpointSecretName, "aws-scim-endpoint-secret-name", "n", config.DefaultAWSSCIMEndpointSecretName, "AWS Secrets Manager secret name for AWS SSO SCIM API Endpoint")
	rootCmd.PersistentFlags().StringVarP(&cfg.AWSS3BucketName, "aws-s3-bucket-name", "b", "", "AWS S3 Bucket name to store the state")
	rootCmd.PersistentFlags().StringVarP(&cfg.AWSS3BucketKey, "aws-s3-bucket-key", "k", config.DefaultAWSS3BucketKey, "AWS S3 Bucket key to store the state")
	rootCmd.PersistentFlags().StringVarP(&cfg.GWSServiceAccountFile, "gws-service-account-file", "s", config.DefaultGWSServiceAccountFile, "Google Workspace service account file")
	rootCmd.PersistentFlags().StringVarP(&cfg.GWSServiceAccountFileSecretName, "gws-service-account-file-secret-name", "o", config.DefaultGWSServiceAccountFileSecretName, "AWS Secrets Manager secret name for Google Workspace service account file")
	rootCmd.PersistentFlags().StringVarP(&cfg.GWSUserEmail, "gws-user-email", "u", "", "GWS user email with allowed access to the Google Workspace Service Account")
	rootCmd.PersistentFlags().StringVarP(&cfg.GWSUserEmailSecretName, "gws-user-email-secret-name", "p", config.DefaultGWSUserEmailSecretName, "AWS Secrets Manager secret name for GWS user email with allowed access to the Google Workspace Service Account")
	rootCmd.Flags().StringSliceVarP(&cfg.GWSGroupsFilter, "gws-groups-filter", "q", []string{""}, "GWS Groups query parameter, example: --gws-groups-filter 'name:Admin* email:admin*' --gws-groups-filter 'name:Power* email:power*'")
	rootCmd.PersistentFlags().StringVarP(&cfg.SyncMethod, "sync-method", "m", config.DefaultSyncMethod, "Sync method to use [groups]")
	rootCmd.PersistentFlags().BoolVarP(&cfg.UseSecretsManager, "use-secrets-manager", "g", config.DefaultUseSecretsManager, "use AWS Secrets Manager content or not (default false)")
}

func run(ctx context.Context) error {
	slog.Debug("viper config", "config", viper.AllSettings())

	if cfg.SyncMethod != "groups" {
		return fmt.Errorf("unknown sync method: %s, only 'groups' are implemented", cfg.SyncMethod)
	}

	slog.Info("starting sync groups", "codeVersion", version.Version)
	timeStart := time.Now()

	ss, err := setup.SyncService(ctx, &cfg)
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
