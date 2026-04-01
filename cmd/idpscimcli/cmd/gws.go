package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

// command gws
var (
	// base gws command
	gwsCmd = &cobra.Command{
		Use:   "gws",
		Short: "Google Workspace commands",
		Long:  `available commands to validate Google Workspace Directory API.`,
	}

	// groups command
	gwsGroupsCmd = &cobra.Command{
		Use:   "groups",
		Short: "Google Workspace Groups commands",
		Long:  `available commands to validate Google Workspace Directory Groups API.`,
	}

	// groups list command
	gwsGroupsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list Groups",
		Long:    `This command is used to list the groups from Google Workspace Directory Service`,
		RunE:    execGWSGroupsList,
	}

	// groups members command
	gwsGroupsMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "Google Workspace Groups Members commands",
		Long:  `available commands to validate Google Workspace Directory Groups Members API.`,
	}

	// groups members list command
	gwsGroupsMembersListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list Members",
		Long:    `This command is used to list the groups members from Google Workspace Directory Service`,
		RunE:    execGWSGroupsMembersList,
	}

	// users command
	gwsUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "Google Workspace Users commands",
		Long:  `available commands to validate Google Workspace Directory Users API.`,
	}

	// user list command
	gwsUsersListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list Users",
		Long:    `This command is used to list the users from Google Workspace Directory Service`,
		RunE:    execGWSUsersList,
	}
)

func init() {
	rootCmd.AddCommand(gwsCmd)

	gwsCmd.AddCommand(gwsGroupsCmd)
	gwsCmd.AddCommand(gwsUsersCmd)

	gwsCmd.PersistentFlags().StringVarP(
		&cfg.GWSServiceAccountFile, "gws-service-account-file", "s",
		config.DefaultGWSServiceAccountFile,
		"path to Google Workspace service account file",
	)

	gwsCmd.PersistentFlags().StringVarP(&cfg.GWSUserEmail,
		"gws-user-email", "u", "",
		"Google Workspace user email with allowed access to the Google Workspace service account",
	)

	// groups command
	gwsGroupsCmd.AddCommand(gwsGroupsListCmd)
	gwsGroupsCmd.AddCommand(gwsGroupsMembersCmd)
	gwsGroupsMembersCmd.AddCommand(gwsGroupsMembersListCmd)

	gwsGroupsListCmd.Flags().StringSliceVarP(
		&cfg.GWSGroupsFilter, "gws-groups-filter", "q", []string{""},
		"GWS Groups query parameter, example: --gws-groups-filter 'name=Admin* email=admin*' --gws-groups-filter 'name=Power* email=power*'",
	)

	gwsGroupsMembersListCmd.Flags().StringSliceVarP(
		&cfg.GWSGroupsFilter, "gws-groups-filter", "q", []string{""},
		"GWS Groups query parameter, example: --gws-groups-filter 'name=Admin* email=admin*' --gws-groups-filter 'name=Power* email=power*'",
	)

	// users command
	gwsUsersCmd.AddCommand(gwsUsersListCmd)
	gwsUsersListCmd.Flags().StringSliceVarP(
		&cfg.GWSUsersFilter, "gws-users-filter", "r", []string{""},
		"GWS Users query parameter, example: --gws-users-filter 'name=Admin* email=admin*' --gws-users-filter 'name=Power* email=power*'",
	)
}

func getGWSDirectoryService(ctx context.Context) (*google.DirectoryService, error) {
	gCreds, err := os.ReadFile(cfg.GWSServiceAccountFile)
	if err != nil {
		return nil, fmt.Errorf("error reading the credentials: %w", err)
	}

	gScopes := []string{
		"https://www.googleapis.com/auth/admin.directory.group.readonly",
		"https://www.googleapis.com/auth/admin.directory.group.member.readonly",
		"https://www.googleapis.com/auth/admin.directory.user.readonly",
	}

	gServiceConfig := google.DirectoryServiceConfig{
		UserEmail:      cfg.GWSUserEmail,
		ServiceAccount: gCreds,
		Scopes:         gScopes,
		Client:         newGWSHTTPClient(),
		UserAgent:      fmt.Sprintf("idp-scim-sync/%s", version.Version),
	}

	gService, err := google.NewService(ctx, gServiceConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating google service: %w", err)
	}

	gDirService, err := google.NewDirectoryService(gService)
	if err != nil {
		return nil, fmt.Errorf("error creating google directory service: %w", err)
	}

	return gDirService, nil
}

func execGWSGroupsList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService, err := getGWSDirectoryService(ctx)
	if err != nil {
		return err
	}

	gGroups, err := gDirService.ListGroups(ctx, cfg.GWSGroupsFilter)
	if err != nil {
		return fmt.Errorf("error listing groups: %w", err)
	}
	slog.Info("groups found", "groups", len(gGroups))

	return show(outFormat, gGroups)
}

func execGWSUsersList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService, err := getGWSDirectoryService(ctx)
	if err != nil {
		return err
	}

	gUsers, err := gDirService.ListUsers(ctx, cfg.GWSUsersFilter)
	if err != nil {
		return fmt.Errorf("error listing users: %w", err)
	}
	slog.Info("users found", "users", len(gUsers))

	return show(outFormat, gUsers)
}

type gwsGroupMembers struct {
	Group   *admin.Group
	Members []*admin.Member
}

func execGWSGroupsMembersList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService, err := getGWSDirectoryService(ctx)
	if err != nil {
		return err
	}

	gGroups, err := gDirService.ListGroups(ctx, cfg.GWSGroupsFilter)
	if err != nil {
		return fmt.Errorf("error listing groups: %w", err)
	}
	slog.Info("groups found", "groups", len(gGroups))

	gMembers := make([]gwsGroupMembers, 0, len(gGroups))
	for _, group := range gGroups {
		members, err := gDirService.ListGroupMembers(ctx, group.Id)
		if err != nil {
			return fmt.Errorf("error listing members for group %s: %w", group.Email, err)
		}
		gMembers = append(gMembers, gwsGroupMembers{
			Group:   group,
			Members: members,
		})
	}

	return show(outFormat, gMembers)
}
