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

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/aws-sso-gws-sync/internal/google"
	"github.com/spf13/cobra"
)

var query []string

// gwsGroupListCmd represents the Google Workspace Groups command
var gwsGroupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list Groups",
	Long:  `This command is used to list the groups from Google Workspace Directory Servive`,
	Run: func(cmd *cobra.Command, args []string) {
		execGWSGroupsList()
	},
}

func init() {
	gwsGroupsCmd.AddCommand(gwsGroupsListCmd)

	gwsGroupsListCmd.Flags().StringSliceVarP(&query, "query", "q", []string{""}, "Google Workspace Groups query parameter, example: --query 'name:Admin* email:admin*' --query 'name:Power* email:power*', see: https://developers.google.com/admin-sdk/directory/v1/guides/search-groups")
}

func execGWSGroupsList() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gCreds, err := ioutil.ReadFile(cfg.ServiceAccountFile)
	if err != nil {
		log.Fatalf("Error reading the credentials", err)
	}

	gDirService, err := google.NewDirectoryService(ctx, cfg.UserEmail, gCreds)
	if err != nil {
		log.Fatalf("Error connecting to google", err)
	}

	gGroups, err := gDirService.ListGroups(query)
	if err != nil {
		log.Fatalf("Error listing groups", err)
	}
	log.Infof("%d groups found", len(gGroups))

	for _, g := range gGroups {
		log.Infof("Name: %s - Email: %s", g.Name, g.Email)
	}
}
