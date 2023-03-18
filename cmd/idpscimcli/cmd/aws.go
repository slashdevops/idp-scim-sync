package cmd

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
	"github.com/spf13/cobra"
)

var filter string

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
		Use:     "config",
		Aliases: []string{"c"},
		Short:   "return Service Provider config",
		Long:    `return AWS SSO SCIM provider config.`,
		RunE:    runAWSServiceConfig,
	}

	// groups command
	awsGroupsCmd = &cobra.Command{
		Use:   "groups",
		Short: "AWS SSO SCIM Groups commands",
		Long:  `available commands to validate AWS SSO SCIM Groups API.`,
	}

	// groups list command
	awsGroupsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "return available groups",
		Long:    `list available groups in AWS SSO SCIM`,
		RunE:    runAWSGroupsList,
	}

	// users command
	awsUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "AWS SSO SCIM Users commands",
		Long:  `available commands to validate AWS SSO SCIM Users API.`,
	}

	// users list command
	awsUsersListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "return available users",
		Long:    `list available usrs in AWS SSO SCIM`,
		RunE:    runAWSUsersList,
	}
)

func init() {
	rootCmd.AddCommand(awsCmd)
	awsCmd.AddCommand(awsGroupsCmd)
	awsCmd.AddCommand(awsUsersCmd)
	awsCmd.AddCommand(awsServiceCmd)
	awsServiceCmd.AddCommand(awsServiceConfigCmd)
	awsGroupsCmd.AddCommand(awsGroupsListCmd)
	awsUsersCmd.AddCommand(awsUsersListCmd)

	awsCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMAccessToken, "aws-scim-access-token", "t", "", "AWS SSO SCIM API Access Token")
	awsCmd.PersistentFlags().StringVarP(&cfg.AWSSCIMEndpoint, "aws-scim-endpoint", "e", "", "AWS SSO SCIM API Endpoint")

	awsGroupsListCmd.Flags().StringVarP(&filter, "filter", "q", "", "AWS SSO SCIM API Filter, example: --filter 'displayName eq \"Group Bar\" and id eq \"12324\"'")
	awsUsersCmd.Flags().StringVarP(&filter, "filter", "q", "", "AWS SSO SCIM API Filter, example: --filter 'displayName eq \"User Bar\" and id eq \"12324\"'")
}

func runAWSServiceConfig(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.MaxIdleConns = 100
	httpTransport.MaxConnsPerHost = 100
	httpTransport.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   maxTimeout,
	}

	awsSCIMService, err := aws.NewSCIMService(httpClient, cfg.AWSSCIMEndpoint, cfg.AWSSCIMAccessToken)
	if err != nil {
		log.Errorf("error creating SCIM service: %s", err.Error())
		return err
	}
	awsSCIMService.UserAgent = "idp-scim-sync/" + version.Version

	awsServiceConfig, err := awsSCIMService.ServiceProviderConfig(ctx)
	if err != nil {
		log.Errorf("error getting service provider config, error: %s", err.Error())
		return err
	}

	show(outFormat, awsServiceConfig)

	return nil
}

func runAWSGroupsList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.MaxIdleConns = 100
	httpTransport.MaxConnsPerHost = 100
	httpTransport.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   maxTimeout,
	}

	awsSCIMService, err := aws.NewSCIMService(httpClient, cfg.AWSSCIMEndpoint, cfg.AWSSCIMAccessToken)
	if err != nil {
		log.Errorf("error creating SCIM service: %s", err.Error())
		return err
	}
	awsSCIMService.UserAgent = "idp-scim-sync/" + version.Version

	awsGroupsResponse, err := awsSCIMService.ListGroups(ctx, filter)
	if err != nil {
		log.Errorf("error listing groups, error: %s", err.Error())
		return err
	}
	log.Infof("%d groups found", awsGroupsResponse.TotalResults)

	show(outFormat, awsGroupsResponse)

	return nil
}

func runAWSUsersList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.MaxIdleConns = 100
	httpTransport.MaxConnsPerHost = 100
	httpTransport.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   maxTimeout,
	}

	awsSCIMService, err := aws.NewSCIMService(httpClient, cfg.AWSSCIMEndpoint, cfg.AWSSCIMAccessToken)
	if err != nil {
		log.Errorf("error creating SCIM service: %s", err.Error())
		return err
	}
	awsSCIMService.UserAgent = "idp-scim-sync/" + version.Version

	awsUsersResponse, err := awsSCIMService.ListUsers(ctx, filter)
	if err != nil {
		log.Errorf("error listing groups, error: %s", err.Error())
		return err
	}
	log.Infof("%d users found", awsUsersResponse.TotalResults)

	show(outFormat, awsUsersResponse)

	return nil
}
