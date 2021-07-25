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

var outFormat string

var awsServiceConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "return Service Provider config",
	Long:  `return AWS SSO SCIM provider config.`,
	Run: func(cmd *cobra.Command, args []string) {
		execAWSServiceConfig()
	},
}

func init() {
	awsServiceCmd.AddCommand(awsServiceConfigCmd)

	awsServiceConfigCmd.Flags().StringVar(&outFormat, "output-format", "json", "output format (json|yaml)")

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
		log.Fatalf("Error creating SCIM service: ", err.Error())
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
