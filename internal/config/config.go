package config

import (
	"encoding/json"
	"log"
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
}

const (
	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = "info"

	// DefaultLogFormat is the default format of the logger
	DefaultLogFormat = "text"

	// DefaultDebug is the default debug status.
	DefaultDebug = false

	// DefaultGoogleCredentials is the default credentials path
	DefaultServiceAccountFile = "credentials.json"
)

// New returns a new Config
func NewConfig() Config {
	return Config{
		Debug:              DefaultDebug,
		LogLevel:           DefaultLogLevel,
		LogFormat:          DefaultLogFormat,
		ServiceAccountFile: DefaultServiceAccountFile,
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
