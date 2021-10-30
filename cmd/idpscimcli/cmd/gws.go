package cmd

import (
	"context"
	"io/ioutil"

	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/pkg/google"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var (
	groupsQuery []string
	usersQuery  []string
)

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
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list Groups",
		Long:    `This command is used to list the groups from Google Workspace Directory Servive`,
		RunE:    execGWSGroupsList,
	}

	// users command
	gwsUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "Google Workspace Users commands",
		Long:  `available commands to validate Google Worspace Directory Users API.`,
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

	gwsCmd.PersistentFlags().StringVarP(&cfg.GWSServiceAccountFile, "gws-service-account-file", "s", config.DefaultGWSServiceAccountFile, "path to Google Workspace service account file")
	gwsCmd.MarkPersistentFlagRequired("gws-service-account-file")

	gwsCmd.PersistentFlags().StringVarP(&cfg.GWSUserEmail, "gws-user-email", "u", "", "Google Workspace user email with allowed access to the Google Workspace service account")
	gwsCmd.MarkPersistentFlagRequired("gws-user-email")

	// groups command
	gwsGroupsCmd.AddCommand(gwsGroupsListCmd)
	gwsGroupsListCmd.Flags().StringSliceVarP(&groupsQuery, "query-groups", "q", []string{""}, "Google Workspace Groups query parameter, example: --query-groups 'name:Admin* email:admin*' --query-groups 'name:Power* email:power*', see: https://developers.google.com/admin-sdk/directory/v1/guides/search-groups")

	// users command
	gwsUsersCmd.AddCommand(gwsUsersListCmd)
	gwsUsersListCmd.Flags().StringSliceVarP(&usersQuery, "query-users", "r", []string{""}, "Google Workspace Users query parameter, example: --query-users 'name:Admin* email:admin*' --query-users 'name:Power* email:power*', see: https://developers.google.com/admin-sdk/directory/v1/guides/search-users")
}

func getGWSDirectoryService(ctx context.Context) *google.DirectoryService {
	gCreds, err := ioutil.ReadFile(cfg.GWSServiceAccountFile)
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

	gGroups, err := gDirService.ListGroups(ctx, groupsQuery)
	if err != nil {
		log.Errorf("error listing groups: %s", err)
		return err
	}
	log.Infof("%d groups found", len(gGroups))

	for _, g := range gGroups {
		log.WithFields(log.Fields{
			"Id":    g.Id,
			"Name":  g.Name,
			"Email": g.Email,
		}).Info("List Group ->")
	}
	return nil
}

func execGWSUsersList(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	gDirService := getGWSDirectoryService(ctx)

	gUsers, err := gDirService.ListUsers(ctx, usersQuery)
	if err != nil {
		log.Errorf("error listing users: %s", err)
		return err
	}
	log.Infof("%d users found", len(gUsers))

	for _, u := range gUsers {
		log.WithFields(log.Fields{
			"Id":        u.Id,
			"Name":      u.Name.FullName,
			"Email":     u.PrimaryEmail,
			"Suspended": u.Suspended,
		}).Info("List User ->")
	}
	return nil
}
