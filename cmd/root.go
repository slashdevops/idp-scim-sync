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
	"github.com/slashdevops/aws-sso-gws-sync/internal/config"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cfg *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-sso-gws-sync",
	Short: "Sync your Google Workspace Groups and Users to AWS Single Sing-On",
	Long: `Keep your AWS Single Sign-On (SSO) users synchronized with your Google Workspace Groups
Sync your Google Workspace Groups and Users to AWS Single Sing-On using
AWS SSO SCIM API (https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html).`,
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cfg = config.NewConfig()

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&cfg.Debug, "debug", "d", config.DefaultDebug, "enable log debug level")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogFormat, "log-format", "f", config.DefaultLogFormat, "set the log format")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogLevel, "log-level", "l", config.DefaultLogLevel, "set the log level")

	rootCmd.PersistentFlags().StringVarP(&cfg.ServiceAccountFile, "gws-service-account-file", "s", config.DefaultServiceAccountFile, "path to Google Workspace service account file")
	rootCmd.PersistentFlags().StringVarP(&cfg.UserEmail, "gws-user-email", "u", "", "Google Workspace user email with allowed access to the Google Workspace Service Account")

	rootCmd.PersistentFlags().StringVarP(&cfg.SCIMAccessToken, "aws-scim-access-token", "t", "", "AWS SSO SCIM API Access Token")
	rootCmd.PersistentFlags().StringVarP(&cfg.SCIMEndpoint, "aws-scim-endpoint", "e", "", "AWS SSO SCIM API Endpoint")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// AWS SSO GWS Sync (asgs)
	viper.SetEnvPrefix("asgs")
	viper.AutomaticEnv()

	envVars := []string{
		"log_format",
		"log_level",
		"gws_service_account_file",
		"gws_user_email",
		"aws_scim_access_token",
		"aws_scim_endpoint",
	}

	for _, e := range envVars {
		if err := viper.BindEnv(e); err != nil {
			log.Fatalf("Cannot bind environment variable", err)
		}
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
}
