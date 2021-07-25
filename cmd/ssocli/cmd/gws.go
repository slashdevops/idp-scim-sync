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
)

var gwsCmd = &cobra.Command{
	Use:   "gws",
	Short: "Google Workspace commands",
	Long:  `available commands to validate Google Worspace Directory API.`,
}

func init() {
	rootCmd.AddCommand(gwsCmd)

	gwsCmd.PersistentFlags().StringVarP(&cfg.ServiceAccountFile, "gws-service-account-file", "s", config.DefaultServiceAccountFile, "path to Google Workspace service account file")
	gwsCmd.MarkPersistentFlagRequired("gws-service-account-file")

	gwsCmd.PersistentFlags().StringVarP(&cfg.UserEmail, "gws-user-email", "u", "", "Google Workspace user email with allowed access to the Google Workspace Service Account")
	gwsCmd.MarkPersistentFlagRequired("gws-user-email")
}
