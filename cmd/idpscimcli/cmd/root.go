package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	cfg        config.Config
	reqTimeout time.Duration
	maxTimeout time.Duration
	outFormat  string
)

// commands root
var rootCmd = &cobra.Command{
	Use:     "idpscimcli",
	Version: version.Version,
	Short:   "Check your  AWS Single Sing-On (SSO) / Google Workspace Groups/Users",
	Long: `
This is a Command-Line Interfaced (CLI) to help you validate and check your source and target Single Sing-On endpoints.
Check your AWS Single Sign-On (SSO) / Google Workspace Groups users and groups and validate your filters over Google Workspace users and groups.`,
}

// Execute validates the configuration and executes the command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cfg = config.New()
	maxTimeout = time.Second * 10

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfg.ConfigFile, "config-file", "c", config.DefaultConfigFile, "configuration file")

	// global configuration for commands
	rootCmd.PersistentFlags().BoolVarP(&cfg.Debug, "debug", "d", config.DefaultDebug, "enable log debug level")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogFormat, "log-format", "f", config.DefaultLogFormat, "set the log format")
	rootCmd.PersistentFlags().StringVarP(&cfg.LogLevel, "log-level", "l", config.DefaultLogLevel, "set the log level")
	rootCmd.PersistentFlags().DurationVarP(&reqTimeout, "timeout", "", maxTimeout, "requests timeout")
	rootCmd.PersistentFlags().StringVar(&outFormat, "output-format", "json", "output format (json|yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("idpscim") // allow to read in from environment

	envVars := []string{
		"log_level",
		"log_format",
		"gws_user_email",
		"gws_service_account_file",
		"gws_groups_filter",
		"gws_users_filter",
		"aws_scim_access_token",
		"aws_scim_endpoint",
	}
	for _, e := range envVars {
		if err := viper.BindEnv(e); err != nil {
			log.Fatalf(errors.Wrap(err, "cannot bind environment variable").Error())
		}
	}

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home)

	currentDir, err := os.Getwd()
	cobra.CheckErr(err)
	viper.AddConfigPath(currentDir)

	fileDir := filepath.Dir(cfg.ConfigFile)
	viper.AddConfigPath(fileDir)

	// Search config in home directory with name "downloader" (without extension).
	fileExtension := filepath.Ext(cfg.ConfigFile)
	fileExtensionName := fileExtension[1:]
	viper.SetConfigType(fileExtensionName)

	fileNameExt := filepath.Base(cfg.ConfigFile)
	fileName := fileNameExt[0 : len(fileNameExt)-len(fileExtension)]
	viper.SetConfigName(fileName)

	log.Debugf("configuration file: dir: %s, name: %s, ext: %s", fileDir, fileName, fileExtension)

	if err := viper.ReadInConfig(); err == nil {
		log.Infof("using config file: %s", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf(errors.Wrap(err, "cannot unmarshal config").Error())
	}

	switch strings.ToLower(cfg.LogFormat) {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Warnf("unknown log format: %s, using text", cfg.LogFormat)
		log.SetFormatter(&log.TextFormatter{})
	}

	if cfg.Debug {
		cfg.LogLevel = "debug"
	}

	// set the configured log level
	if level, err := log.ParseLevel(strings.ToLower(cfg.LogLevel)); err == nil {
		log.SetLevel(level)
	}
}
