/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"
	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ssolambda",
	Version: version.Version,
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	cfg = config.NewConfig()
	cfg.IsLambda = len(os.Getenv("_LAMBDA_SERVER_PORT")) > 0

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&cfg.Debug, "debug", "d", config.DefaultDebug, "enable log debug level")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogFormat, "log-format", "f", config.DefaultLogFormat, "set the log format")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogLevel, "log-level", "l", config.DefaultLogLevel, "set the log level")

	rootCmd.PersistentFlags().StringVarP(&cfg.SCIMAccessToken, "aws-scim-access-token", "t", "", "AWS SSO SCIM API Access Token")
	rootCmd.MarkPersistentFlagRequired("aws-scim-access-token")

	rootCmd.PersistentFlags().StringVarP(&cfg.SCIMEndpoint, "aws-scim-endpoint", "e", "", "AWS SSO SCIM API Endpoint")
	rootCmd.MarkPersistentFlagRequired("aws-scim-endpoint")

	rootCmd.PersistentFlags().StringVarP(&cfg.ServiceAccountFile, "gws-service-account-file", "s", config.DefaultServiceAccountFile, "path to Google Workspace service account file")
	rootCmd.MarkPersistentFlagRequired("gws-service-account-file")

	rootCmd.PersistentFlags().StringVarP(&cfg.UserEmail, "gws-user-email", "u", "", "Google Workspace user email with allowed access to the Google Workspace Service Account")
	rootCmd.MarkPersistentFlagRequired("gws-user-email")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	switch cfg.LogFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Fatal("Unsupported log format")
	}

	if cfg.Debug {
		cfg.LogLevel = "debug"
	}

	// set the configured log level
	if level, err := log.ParseLevel(cfg.LogLevel); err == nil {
		log.SetLevel(level)
	}

	if cfg.IsLambda {
		configLambda()
	}
}

func configLambda() {
	awsconf, err := awsconf.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf(errors.Wrap(err, "cannot load aws config").Error())
	}

	svc := secretsmanager.NewFromConfig(awsconf)

	secrets, err := aws.NewSecretsManagerService(svc)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "cannot create secrets manager service").Error())
	}

	unwrap, err := secrets.GetSecretValue(context.TODO(), "SSOLambdaGoogleUserEmail")
	if err != nil {
		log.Fatalf(errors.Wrap(err, "cannot get secretmanager value").Error())
	}
	cfg.UserEmail = unwrap

	unwrap, err = secrets.GetSecretValue(context.TODO(), "SSOLambdaGoogleCredentialsFile")
	if err != nil {
		log.Fatalf(errors.Wrap(err, "cannot get secretmanager value").Error())
	}
	cfg.ServiceAccountFile = unwrap

	unwrap, err = secrets.GetSecretValue(context.TODO(), "SCIMAccessToken")
	if err != nil {
		log.Fatalf(errors.Wrap(err, "cannot get secretmanager value").Error())
	}
	cfg.SCIMAccessToken = unwrap

	unwrap, err = secrets.GetSecretValue(context.TODO(), "SCIMEndpoint")
	if err != nil {
		log.Fatalf(errors.Wrap(err, "cannot get secretmanager value").Error())
	}
	cfg.SCIMEndpoint = unwrap
}
