package cmd

import (
	"context"
	"log/slog"

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
		Long:    `list available users in AWS SSO SCIM`,
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

// newAWSSCIMService creates an AWS SCIM service with the configured HTTP client.
func newAWSSCIMService() (*aws.SCIMService, error) {
	svc, err := aws.NewSCIMService(newSCIMHTTPClient(), cfg.AWSSCIMEndpoint, cfg.AWSSCIMAccessToken)
	if err != nil {
		return nil, err
	}
	svc.UserAgent = "idp-scim-sync/" + version.Version
	return svc, nil
}

func runAWSServiceConfig(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	svc, err := newAWSSCIMService()
	if err != nil {
		return err
	}

	result, err := svc.ServiceProviderConfig(ctx)
	if err != nil {
		return err
	}

	return show(outFormat, result)
}

func runAWSGroupsList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	svc, err := newAWSSCIMService()
	if err != nil {
		return err
	}

	result, err := svc.ListGroups(ctx, filter)
	if err != nil {
		return err
	}
	slog.Info("groups found", "groups", result.TotalResults)

	return show(outFormat, result)
}

func runAWSUsersList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	svc, err := newAWSSCIMService()
	if err != nil {
		return err
	}

	result, err := svc.ListUsers(ctx, filter)
	if err != nil {
		return err
	}
	slog.Info("users found", "users", result.TotalResults)

	return show(outFormat, result)
}
