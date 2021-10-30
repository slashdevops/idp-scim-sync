package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	log "github.com/sirupsen/logrus"
)

var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "idpscim",
	Version: version.Version,
	Short:   "Sync your AWS Single Sing-On (SSO) with Google Workspace",
	Long: `
Sync your Google Workspace Groups and Users to AWS Single Sing-On using
AWS SSO SCIM API (https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html).`,
	RunE: sync,
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
	cfg.IsLambda = len(os.Getenv("_LAMBDA_SERVER_PORT")) > 0

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&cfg.Debug, "debug", "d", config.DefaultDebug, "enable log debug level")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogFormat, "log-format", "f", config.DefaultLogFormat, "set the log format")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogLevel, "log-level", "l", config.DefaultLogLevel, "set the log level")

	rootCmd.PersistentFlags().StringVarP(&cfg.SCIMAccessToken, "aws-scim-access-token", "t", "", "AWS SSO SCIM API Access Token")
	_ = rootCmd.MarkPersistentFlagRequired("aws-scim-access-token")

	rootCmd.PersistentFlags().StringVarP(&cfg.SCIMEndpoint, "aws-scim-endpoint", "e", "", "AWS SSO SCIM API Endpoint")
	_ = rootCmd.MarkPersistentFlagRequired("aws-scim-endpoint")

	rootCmd.PersistentFlags().StringVarP(&cfg.GWSServiceAccountFile, "gws-service-account-file", "s", config.DefaultGWSServiceAccountFile, "path to Google Workspace service account file")
	_ = rootCmd.MarkPersistentFlagRequired("gws-service-account-file")

	rootCmd.PersistentFlags().StringVarP(&cfg.GWSUserEmail, "gws-user-email", "u", "", "Google Workspace user email with allowed access to the Google Workspace Service Account")
	_ = rootCmd.MarkPersistentFlagRequired("gws-user-email")

	rootCmd.Flags().StringSliceVarP(&cfg.GWSGroupsFilter, "query-groups", "q", []string{""}, "Google Workspace Groups query parameter, example: --query-groups 'name:Admin* email:admin*' --query-groups 'name:Power* email:power*', see: https://developers.google.com/admin-sdk/directory/v1/guides/search-groups")
	rootCmd.Flags().StringSliceVarP(&cfg.GWSUsersFilter, "query-users", "r", []string{""}, "Google Workspace Users query parameter, example: --query-users 'name:Admin* email:admin*' --query-users 'name:Power* email:power*', see: https://developers.google.com/admin-sdk/directory/v1/guides/search-users")

	rootCmd.PersistentFlags().StringVarP(&cfg.SyncMethod, "sync-method", "m", config.DefaultSyncMethod, "Sync method to use (groups)")

	rootCmd.PersistentFlags().StringVarP(&cfg.AWSS3BucketName, "aws-s3-bucket-name", "b", "", "AWS S3 Bucket name to store the state")
	_ = rootCmd.MarkPersistentFlagRequired("aws-s3-bucket-name")

	rootCmd.PersistentFlags().StringVarP(&cfg.AWSS3BucketKey, "aws-s3-bucket-key", "k", "", "AWS S3 Bucket key to store the state")
	_ = rootCmd.MarkPersistentFlagRequired("aws-s3-bucket-key")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("idpscim") // allow to read in from environment
	viper.AutomaticEnv()          // read in environment variables that match

	envVars := []string{
		"log_level",
		"log_format",
		"sync_method",
		"aws_s3_bucket_name",
		"aws_s3_bucket_key",
		"gws_user_email",
		"gws_service_account_file",
		"gws_service_account_file_secret_name",
		"gws_user_email_secret_name",
		"gws_groups_filter",
		"gws_users_filter",
		"scim_access_token",
		"scim_endpoint",
		"scim_endpoint_secret_name",
		"scim_access_token_secret_name",
	}

	for _, e := range envVars {
		if err := viper.BindEnv(e); err != nil {
			log.Fatalf(errors.Wrap(err, "idpscim: cannot bind environment variable").Error())
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "idpscim: using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot unmarshal config").Error())
	}

	switch cfg.LogFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Warnf("idpscim: unknown log format: %s, using text", cfg.LogFormat)
		log.SetFormatter(&log.TextFormatter{})
	}

	if cfg.Debug {
		cfg.LogLevel = "debug"
	}

	// set the configured log level
	if level, err := log.ParseLevel(cfg.LogLevel); err == nil {
		log.SetLevel(level)
	}

	if cfg.IsLambda {
		getSecrets()
	}
}

func getSecrets() {
	awsconf, err := awsconf.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot load aws config").Error())
	}

	svc := secretsmanager.NewFromConfig(awsconf)

	secrets, err := aws.NewSecretsManagerService(svc)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot create aws secrets manager service").Error())
	}

	unwrap, err := secrets.GetSecretValue(context.Background(), cfg.GWSUserEmailSecretName)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot get secretmanager value").Error())
	}
	cfg.GWSUserEmail = unwrap

	unwrap, err = secrets.GetSecretValue(context.Background(), cfg.GWSServiceAccountFileSecretName)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot get secretmanager value").Error())
	}
	cfg.GWSServiceAccountFile = unwrap

	unwrap, err = secrets.GetSecretValue(context.Background(), cfg.SCIMAccessTokenSecretName)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot get secretmanager value").Error())
	}
	cfg.SCIMAccessToken = unwrap

	unwrap, err = secrets.GetSecretValue(context.Background(), cfg.SCIMEndpointSecretName)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot get secretmanager value").Error())
	}
	cfg.SCIMEndpoint = unwrap
}

func sync(cmd *cobra.Command, args []string) error {
	if cfg.SyncMethod == "groups" {
		return syncGroups(cmd, args)
	}
	return nil
}

func syncGroups(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	serviceAccount, err := ioutil.ReadFile(cfg.GWSServiceAccountFile)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot read service account file").Error())
	}

	// Google Client Service
	gwsService, err := google.NewService(ctx, cfg.GWSUserEmail, serviceAccount)
	if err != nil {
		return errors.Wrap(err, "idpscim: cannot create google service")
	}

	// Google Directory Service
	gds, err := google.NewDirectoryService(gwsService)
	if err != nil {
		return errors.Wrap(err, "idpscim: cannot create google directory service")
	}

	// Identity Provider Service
	idp, err := idp.NewIdentityProvider(gds)
	if err != nil {
		return errors.Wrap(err, "idpscim: cannot create identity provider service")
	}

	// httpClient
	httpClient := &http.Client{}

	// AWS SCIM Service
	awsSCIM, err := aws.NewSCIMService(httpClient, cfg.SCIMEndpoint, cfg.SCIMAccessToken)
	if err != nil {
		return errors.Wrap(err, "idpscim: cannot create aws scim service")
	}

	scim, err := scim.NewSCIMProvider(awsSCIM)
	if err != nil {
		return errors.Wrap(err, "idpscim: cannot create scim provider")
	}

	awsconf, err := awsconf.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot load aws config").Error())
	}

	s3Client := s3.NewFromConfig(awsconf)
	repo, err := repository.NewS3Repository(s3Client, repository.WithBucket(cfg.AWSS3BucketName), repository.WithKey(cfg.AWSS3BucketKey))
	if err != nil {
		log.Fatalf(errors.Wrap(err, "idpscim: cannot create s3 repository").Error())
	}

	ss, err := core.NewSyncService(ctx, idp, scim, repo, core.WithIdentityProviderGroupsFilter(cfg.GWSGroupsFilter))
	if err != nil {
		return errors.Wrap(err, "idpscim: cannot create sync service")
	}

	if err := ss.SyncGroupsAndTheirMembers(); err != nil {
		return errors.Wrap(err, "idpscim: cannot sync groups and their members")
	}

	return nil
}
