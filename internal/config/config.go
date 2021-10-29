package config

import (
	"encoding/json"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	IsLambda bool
	Debug    bool

	LogLevel  string `mapstructure:"log_level"`
	LogFormat string `mapstructure:"log_format"`

	ServiceAccountFile string `mapstructure:"service_account_file"`
	UserEmail          string `mapstructure:"user_email"`

	SCIMEndpoint    string `mapstructure:"scim_endpoint"`
	SCIMAccessToken string `mapstructure:"scim_access_token"`

	// SyncMethod allow to defined the sync method used to get the user and groups from Google Workspace
	SyncMethod string `mapstructure:"sync_method"`
}

const (
	// DefaultIsLambda is the progam execute as a lambda function?
	DefaultIsLambda = false

	// DefaultLogLevel is the default logging level.
	// possible values: "debug", "info", "warn", "error", "fatal", "panic"
	DefaultLogLevel = "info"

	// DefaultLogFormat is the default format of the logger
	// possible values: "text", "json"
	DefaultLogFormat = "text"

	// DefaultDebug is the default debug status.
	DefaultDebug = false

	// DefaultGoogleCredentials is the default credentials path
	DefaultServiceAccountFile = "credentials.json"

	// DefaultSyncMethod is the default sync method to use.
	DefaultSyncMethod = "groups"
)

// New returns a new Config
func New() Config {
	return Config{
		IsLambda:           DefaultIsLambda,
		Debug:              DefaultDebug,
		LogLevel:           DefaultLogLevel,
		LogFormat:          DefaultLogFormat,
		ServiceAccountFile: DefaultServiceAccountFile,
		SyncMethod:         DefaultSyncMethod,
	}
}

// toJSON return a json pretty of the config.
func (c *Config) toJSON() []byte {
	JSON, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return JSON
}

// toYAML return a yaml of the config.
func (c *Config) toYAML() []byte {
	YAML, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return YAML
}
