package cmd

import (
	"time"

	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var (
	cfg        config.Config
	reqTimeout time.Duration
)

// commands root
var rootCmd = &cobra.Command{
	Use:     "idpscimcli",
	Version: version.Version,
	Short:   "Check your  AWS Single Sing-On (SSO) / Google Workspace Groups/Users",
	Long: `
This is a Command-Line Interfaced (CLI) to help you validate and check your source and target Single Sing-On endpoints.
Check your AWS Single Sign-On (SSO) / Google Workspace Groups users and groups and validate your filters over Google Worspace users and groups.`,
}

// Execute validates the configuration and executes the command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cfg = config.New()

	cobra.OnInitialize(initConfig)

	// global configuration for commands
	rootCmd.PersistentFlags().BoolVarP(&cfg.Debug, "debug", "d", config.DefaultDebug, "enable log debug level")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogFormat, "log-format", "f", config.DefaultLogFormat, "set the log format")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogLevel, "log-level", "l", config.DefaultLogLevel, "set the log level")
	rootCmd.PersistentFlags().DurationVarP(&reqTimeout, "timeout", "", time.Second*10, "requests timeout")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	switch cfg.LogFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Fatal("Unsupported log format")
	}

	if cfg.Debug {
		cfg.LogLevel = "debug"
	}

	// set the configured log level
	if level, err := log.ParseLevel(cfg.LogLevel); err == nil {
		log.SetLevel(level)
	}
}
