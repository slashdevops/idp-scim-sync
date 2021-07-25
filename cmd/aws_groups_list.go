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
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/aws-sso-gws-sync/internal/aws"
	"github.com/spf13/cobra"
)

var filter []string

var awsGroupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "return available groups",
	Long:  `list available groups in AWS SSO`,
	Run: func(cmd *cobra.Command, args []string) {
		execAWSGroupsList()
	},
}

func init() {
	awsGroupsCmd.AddCommand(awsGroupsListCmd)

	awsGroupsListCmd.Flags().StringSliceVarP(&filter, "filter", "q", []string{""}, "AWS SSO SCIM API Filter, example: --filter 'displayName eq \"Group Bar\" and id eq 12324' --filter 'displayName eq \"Group Foo\"' see: https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html#examples-filter-listgroups")
}

func execAWSGroupsList() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.MaxIdleConns = 100
	httpTransport.MaxConnsPerHost = 100
	httpTransport.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   time.Second * 10,
	}

	awsSCIMService, err := aws.NewSCIMService(&ctx, httpClient, cfg.SCIMEndpoint, cfg.SCIMAccessToken)
	if err != nil {
		log.Fatalf("Error creating SCIM service: ", err.Error())
	}

	awsGroups, err := awsSCIMService.ListGroups(filter)
	if err != nil {
		log.Fatalf("Error listing groups, error: %s", err.Error())
	}
	log.Infof("%d groups found", len(awsGroups))

	for _, g := range awsGroups {
		log.Infof("Name: %s - ID: %s", g.DisplayName, g.ID)
	}
}
