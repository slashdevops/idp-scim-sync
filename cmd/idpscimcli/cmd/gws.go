package cmd

import (
	"context"
	"os"

	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"

	log "github.com/sirupsen/logrus"
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
		"GWS Groups query parameter, example: --gws-groups-filter 'name:Admin* email:admin*' --gws-groups-filter 'name:Power* email:power*'",
	)

	gwsGroupsMembersListCmd.Flags().StringSliceVarP(
		&cfg.GWSGroupsFilter, "gws-groups-filter", "q", []string{""},
		"GWS Groups query parameter, example: --gws-groups-filter 'name:Admin* email:admin*' --gws-groups-filter 'name:Power* email:power*'",
	)

	// users command
	gwsUsersCmd.AddCommand(gwsUsersListCmd)
	gwsUsersListCmd.Flags().StringSliceVarP(
		&cfg.GWSUsersFilter, "gws-users-filter", "r", []string{""},
		"GWS Users query parameter, example: --gws-users-filter 'name:Admin* email:admin*' --gws-users-filter 'name:Power* email:power*'",
	)
}

func getGWSDirectoryService(ctx context.Context) *google.DirectoryService {
	gCreds, err := os.ReadFile(cfg.GWSServiceAccountFile)
	if err != nil {
		log.Fatalf("error reading the credentials: %s", err)
	}

	gScopes := []string{
		"https://www.googleapis.com/auth/admin.directory.group.readonly",
		"https://www.googleapis.com/auth/admin.directory.group.member.readonly",
		"https://www.googleapis.com/auth/admin.directory.user.readonly",
	}

	gService, err := google.NewService(ctx, cfg.GWSUserEmail, gCreds, gScopes...)
	if err != nil {
		log.Fatalf("error creating service: %s", err)
	}

	gDirService, err := google.NewDirectoryService(gService)
	if err != nil {
		log.Fatalf("error creating directory service: %s", err)
	}

	return gDirService
}

func execGWSGroupsList(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService := getGWSDirectoryService(ctx)

	gGroups, err := gDirService.ListGroups(ctx, cfg.GWSGroupsFilter)
	if err != nil {
		log.Errorf("error listing groups: %s", err)
		return err
	}
	log.Infof("%d groups found", len(gGroups))

	show(outFormat, gGroups)

	return nil
}

func execGWSUsersList(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService := getGWSDirectoryService(ctx)

	gUsers, err := gDirService.ListUsers(ctx, cfg.GWSUsersFilter)
	if err != nil {
		log.Errorf("error listing users: %s", err)
		return err
	}
	log.Infof("%d users found", len(gUsers))

	show(outFormat, gUsers)

	return nil
}

func execGWSGroupsMembersList(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService := getGWSDirectoryService(ctx)

	gGroups, err := gDirService.ListGroups(ctx, cfg.GWSGroupsFilter)
	if err != nil {
		log.Errorf("error listing groups: %s", err)
		return err
	}

	log.Infof("%d groups found", len(gGroups))
	var gMembers []*admin.Member

	for _, group := range gGroups {
		members, err := gDirService.ListGroupMembers(ctx, group.Id)
		if err != nil {
			log.Errorf("error listing group members: %s", err)
			return err
		}
		gMembers = append(gMembers, members...)
	}

	show(outFormat, gMembers)

	return nil
}
