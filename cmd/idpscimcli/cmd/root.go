// Package cmd provides the root command and configuration for the idpscimcli CLI application.
package cmd

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/slashdevops/httpx"
	"github.com/slashdevops/idp-scim-sync/internal/config"
	"github.com/slashdevops/idp-scim-sync/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Short:   "Check your AWS Single Sign-On (SSO) / Google Workspace Groups/Users",
	Long: `
This is a Command-Line Interfaced (CLI) to help you validate and check your source and target Single Sign-On endpoints.
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
			slog.Error("cannot bind environment variable", "error", err)
			os.Exit(1)
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

	fileExtension := filepath.Ext(cfg.ConfigFile)
	fileExtensionName := fileExtension[1:]
	viper.SetConfigType(fileExtensionName)

	fileNameExt := filepath.Base(cfg.ConfigFile)
	fileName := strings.TrimSuffix(fileNameExt, fileExtension)
	viper.SetConfigName(fileName)

	slog.Debug("configuration file", "dir", fileDir, "name", fileName, "ext", fileExtension)

	if err := viper.ReadInConfig(); err == nil {
		slog.Info("using config file", "file", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("cannot unmarshal config", "error", err)
		os.Exit(1)
	}

	if cfg.Debug {
		cfg.LogLevel = "debug"
	}

	// Configure logger after config is parsed so log level and format are known
	var logHandlerOptions *slog.HandlerOptions
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}
	case "info":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelInfo}
	case "warn":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelWarn}
	case "error":
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelError, AddSource: true}
	default:
		slog.Warn("unknown log level, setting it to info", "level", cfg.LogLevel)
		logHandlerOptions = &slog.HandlerOptions{Level: slog.LevelInfo}
	}

	var logHandler slog.Handler
	switch strings.ToLower(cfg.LogFormat) {
	case "json":
		logHandler = slog.NewJSONHandler(os.Stdout, logHandlerOptions)
	case "text":
		logHandler = slog.NewTextHandler(os.Stdout, logHandlerOptions)
	default:
		slog.Warn("unknown log format, using text", "format", cfg.LogFormat)
		logHandler = slog.NewTextHandler(os.Stdout, logHandlerOptions)
	}

	slog.SetDefault(slog.New(logHandler))
}

// newSCIMHTTPClient creates an HTTP client configured for AWS SCIM API calls
// with jitter backoff and connection pooling.
func newSCIMHTTPClient() *http.Client {
	return httpx.NewClientBuilder().
		WithMaxRetries(10).
		WithRetryStrategy(httpx.JitterBackoffStrategy).
		WithRetryBaseDelay(500 * time.Millisecond).
		WithRetryMaxDelay(10 * time.Second).
		WithMaxIdleConns(100).
		WithMaxIdleConnsPerHost(100).
		Build()
}

// newGWSHTTPClient creates an HTTP client configured for Google Workspace API calls
// with exponential backoff.
func newGWSHTTPClient() *http.Client {
	return httpx.NewClientBuilder().
		WithMaxRetries(3).
		WithRetryStrategy(httpx.ExponentialBackoffStrategy).
		WithRetryBaseDelay(500 * time.Millisecond).
		WithRetryMaxDelay(10 * time.Second).
		Build()
}
