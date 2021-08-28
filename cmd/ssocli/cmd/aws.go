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
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"github.com/spf13/cobra"
)

var (
	outFormat       string
	filter          string
	SCIMEndpoint    string
	SCIMAccessToken string
)

// commands aws
var (
	// base aws command
	awsCmd = &cobra.Command{
		Use:   "aws",
		Short: "AWS SSO SCIM commands",
		Long:  `Available commands for AWS SSO SCIM.`,
	}

	// service command
	awsServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "AWS SSO SCIM Service commands",
		Long:  `available commands to validate AWS SSO SCIM Service API.`,
	}

	// service config command
	awsServiceConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "return Service Provider config",
		Long:  `return AWS SSO SCIM provider config.`,
		Run: func(cmd *cobra.Command, args []string) {
			execAWSServiceConfig()
		},
	}

	// groups command
	awsGroupsCmd = &cobra.Command{
		Use:   "groups",
		Short: "AWS SSO SCIM Groups commands",
		Long:  `available commands to validate AWS SSO SCIM Groups API.`,
	}

	// groups list command
	awsGroupsListCmd = &cobra.Command{
		Use:   "list",
		Short: "return available groups",
		Long:  `list available groups in AWS SSO`,
		Run: func(cmd *cobra.Command, args []string) {
			execAWSGroupsList()
		},
	}
)

func init() {
	rootCmd.AddCommand(awsCmd)
	awsCmd.AddCommand(awsGroupsCmd)

	awsCmd.PersistentFlags().StringVarP(&cfg.SCIMAccessToken, "aws-scim-access-token", "t", "", "AWS SSO SCIM API Access Token")
	awsCmd.MarkPersistentFlagRequired("aws-scim-access-token")

	awsCmd.PersistentFlags().StringVarP(&cfg.SCIMEndpoint, "aws-scim-endpoint", "e", "", "AWS SSO SCIM API Endpoint")
	awsCmd.MarkPersistentFlagRequired("aws-scim-endpoint")

	awsCmd.AddCommand(awsServiceCmd)
	awsServiceCmd.AddCommand(awsServiceConfigCmd)
	awsServiceConfigCmd.Flags().StringVar(&outFormat, "output-format", "json", "output format (json|yaml)")

	awsGroupsCmd.AddCommand(awsGroupsListCmd)

	awsGroupsListCmd.Flags().StringVarP(&filter, "filter", "q", "", "AWS SSO SCIM API Filter, example: --filter 'displayName eq \"Group Bar\" and id eq \"12324\"', see: https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html#examples-filter-listgroups")
	awsGroupsListCmd.Flags().StringVar(&outFormat, "output-format", "json", "output format (json|yaml)")
}

func execAWSServiceConfig() {
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
		log.Fatalf("Error creating SCIM service: %s", err.Error())
	}

	awsServiceConfig, err := awsSCIMService.ServiceProviderConfig()
	if err != nil {
		log.Fatalf("Error getting service provider config, error: %s", err.Error())
	}

	switch outFormat {
	case "json":
		log.Printf("%s", awsServiceConfig.ToJSON())
	case "yaml":
		log.Printf("%s", awsServiceConfig.ToYAML())
	default:
		log.Fatalf("Unknown output format: %s", outFormat)
	}
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
		log.Fatalf("Error creating SCIM service: %s", err.Error())
	}

	awsGroupsResponse, err := awsSCIMService.ListGroups(filter)
	if err != nil {
		log.Fatalf("Error listing groups, error: %s", err.Error())
	}
	log.Infof("%d groups found", awsGroupsResponse.TotalResults)

	switch outFormat {
	case "json":
		log.Printf("%s", awsGroupsResponse.ToJSON())
	case "yaml":
		log.Printf("%s", awsGroupsResponse.ToYAML())
	default:
		log.Fatalf("Unknown output format: %s", outFormat)
	}
}
