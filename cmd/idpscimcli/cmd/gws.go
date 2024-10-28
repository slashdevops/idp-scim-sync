package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/slashdevops/idp-scim-sync/internal/config"
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
		Long:    `This command is used to list the groups from Google Workspace Directory Servive`,
		RunE:    execGWSGroupsList,
	}

	// groups members command
	gwsGroupsMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "Google Workspace Groups Members commands",
		Long:  `available commands to validate Google Workspace Directory Groups Members API.`,
	}

	// groups list command
	gwsGroupsMembersListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list Members",
		Long:    `This command is used to list the groups members from Google Workspace Directory Servive`,
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
		Long:    `This command is used to list the users from Google Workspace Directory Servive`,
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

func getGWSDirectoryService(ctx context.Context) *google.DirectoryService {
	gCreds, err := os.ReadFile(cfg.GWSServiceAccountFile)
	if err != nil {
		slog.Error("error reading the credentials", "error", err)
		os.Exit(1)
	}

	gScopes := []string{
		"https://www.googleapis.com/auth/admin.directory.group.readonly",
		"https://www.googleapis.com/auth/admin.directory.group.member.readonly",
		"https://www.googleapis.com/auth/admin.directory.user.readonly",
	}

	gService, err := google.NewService(ctx, cfg.GWSUserEmail, gCreds, gScopes...)
	if err != nil {
		slog.Error("error creating service", "error", err)
		os.Exit(1)
	}

	gDirService, err := google.NewDirectoryService(gService)
	if err != nil {
		slog.Error("error creating directory service", "error", err)
		os.Exit(1)
	}

	return gDirService
}

func execGWSGroupsList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService := getGWSDirectoryService(ctx)

	gGroups, err := gDirService.ListGroups(ctx, cfg.GWSGroupsFilter)
	if err != nil {
		slog.Error("error listing groups", "error", err)
		return err
	}
	slog.Info("groups found", "groups", len(gGroups))

	show(outFormat, gGroups)

	return nil
}

func execGWSUsersList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService := getGWSDirectoryService(ctx)

	gUsers, err := gDirService.ListUsers(ctx, cfg.GWSUsersFilter)
	if err != nil {
		slog.Error("error listing users", "error", err)
		return err
	}
	slog.Info("users found", "users", len(gUsers))

	show(outFormat, gUsers)

	return nil
}

type gwsGroupMembers struct {
	Group   *admin.Group
	Members []*admin.Member
}

func execGWSGroupsMembersList(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService := getGWSDirectoryService(ctx)

	gGroups, err := gDirService.ListGroups(ctx, cfg.GWSGroupsFilter)
	if err != nil {
		slog.Error("error listing groups", "error", err)
		return err
	}
	slog.Info("groups found", "groups", len(gGroups))

	gMembers := make([]gwsGroupMembers, 0)

	for _, group := range gGroups {
		members, err := gDirService.ListGroupMembers(ctx, group.Id)
		if err != nil {
			slog.Error("error listing group members", "error", err)
			return err
		}
		e := gwsGroupMembers{
			Group:   group,
			Members: members,
		}
		gMembers = append(gMembers, e)
	}

	show(outFormat, gMembers)

	return nil
}
