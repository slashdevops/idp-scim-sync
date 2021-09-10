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
	"io/ioutil"

	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var query []string

// command gws
var (
	// base gws command
	gwsCmd = &cobra.Command{
		Use:   "gws",
		Short: "Google Workspace commands",
		Long:  `available commands to validate Google Worspace Directory API.`,
	}

	// groups command
	gwsGroupsCmd = &cobra.Command{
		Use:   "groups",
		Short: "Google Workspace Groups commands",
		Long:  `available commands to validate Google Worspace Directory Groups API.`,
	}

	// groups list command
	gwsGroupsListCmd = &cobra.Command{
		Use:   "list",
		Short: "list Groups",
		Long:  `This command is used to list the groups from Google Workspace Directory Servive`,
		Run: func(cmd *cobra.Command, args []string) {
			execGWSGroupsList()
		},
	}
)

func init() {
	rootCmd.AddCommand(gwsCmd)

	gwsCmd.AddCommand(gwsGroupsCmd)

	gwsCmd.PersistentFlags().StringVarP(&cfg.ServiceAccountFile, "gws-service-account-file", "s", config.DefaultServiceAccountFile, "path to Google Workspace service account file")
	gwsCmd.MarkPersistentFlagRequired("gws-service-account-file")

	gwsCmd.PersistentFlags().StringVarP(&cfg.UserEmail, "gws-user-email", "u", "", "Google Workspace user email with allowed access to the Google Workspace Service Account")
	gwsCmd.MarkPersistentFlagRequired("gws-user-email")

	gwsGroupsCmd.AddCommand(gwsGroupsListCmd)
	gwsGroupsListCmd.Flags().StringSliceVarP(&query, "query", "q", []string{""}, "Google Workspace Groups query parameter, example: --query 'name:Admin* email:admin*' --query 'name:Power* email:power*', see: https://developers.google.com/admin-sdk/directory/v1/guides/search-groups")
}

func execGWSGroupsList() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gCreds, err := ioutil.ReadFile(cfg.ServiceAccountFile)
	if err != nil {
		log.Fatalf("Error reading the credentials: %s", err)
	}

	gScopes := []string{
		"https://www.googleapis.com/auth/admin.directory.group.readonly",
		"https://www.googleapis.com/auth/admin.directory.group.member.readonly",
		"https://www.googleapis.com/auth/admin.directory.user.readonly",
	}

	gService, err := google.NewService(ctx, cfg.UserEmail, gCreds, gScopes...)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	gDirService, err := google.NewDirectoryService(ctx, gService)
	if err != nil {
		log.Fatalf("Error creating directory service: %s", err)
	}

	gGroups, err := gDirService.ListGroups(query)
	if err != nil {
		log.Fatalf("Error listing groups: %s", err)
	}
	log.Infof("%d groups found", len(gGroups))

	for _, g := range gGroups {
		log.Infof("Name: %s - Email: %s", g.Name, g.Email)
	}
}
