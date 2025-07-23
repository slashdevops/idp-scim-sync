// Package config provides the configuration for the application.
package config

import "fmt"

const (
	// DefaultIsLambda is the program execute as a lambda function?
	DefaultIsLambda = false

	// DefaultLogLevel is the default logging level.
	// possible values: "debug", "info", "warn", "error", "fatal", "panic"
	DefaultLogLevel = "info"

	// DefaultLogFormat is the default format of the logger
	// possible values: "text", "json"
	DefaultLogFormat = "text"

	// DefaultDebug is the default debug status.
	DefaultDebug = false

	// DefaultGWSServiceAccountFile is the name of the file containing the service account credentials.
	DefaultGWSServiceAccountFile = "credentials.json"

	// DefaultSyncMethod is the default sync method to use.
	DefaultSyncMethod = "groups"

	// DefaultAWSS3BucketKey is the key of the AWS S3 bucket.
	DefaultAWSS3BucketKey = "state.json"

	// DefaultConfigFile is the default config file name.
	DefaultConfigFile = ".idpscim.yaml"

	// DefaultGWSServiceAccountFileSecretName is the name of the secret containing the service account credentials.
	DefaultGWSServiceAccountFileSecretName = "IDPSCIM_GWSServiceAccountFile"

	// DefaultGWSUserEmailSecretName is the name of the secret containing the user email.
	DefaultGWSUserEmailSecretName = "IDPSCIM_GWSUserEmail"

	// DefaultAWSSCIMEndpointSecretName is the name of the secret containing the SCIM endpoint.
	DefaultAWSSCIMEndpointSecretName = "IDPSCIM_SCIMEndpoint"

	// DefaultAWSSCIMAccessTokenSecretName is the name of the secret containing the SCIM access token.
	DefaultAWSSCIMAccessTokenSecretName = "IDPSCIM_SCIMAccessToken"

	// DefaultUseSecretsManager determines if we will use the AWS Secrets Manager secrets or program parameter values
	DefaultUseSecretsManager = false
)

var (
	// ErrInvalidLogLevel is returned when the log level is invalid.
	ErrInvalidLogLevel = fmt.Errorf("invalid log level")
	// ErrInvalidLogFormat is returned when the log format is invalid.
	ErrInvalidLogFormat = fmt.Errorf("invalid log format")
	// ErrMissingAWSSCIMEndpoint is returned when the AWS SCIM endpoint is missing.
	ErrMissingAWSSCIMEndpoint = fmt.Errorf("missing AWS SCIM endpoint")
	// ErrMissingAWSSCIMAccessToken is returned when the AWS SCIM access token is missing.
	ErrMissingAWSSCIMAccessToken = fmt.Errorf("missing AWS SCIM access token")
	// ErrMissingGWSServiceAccountFile is returned when the GWS service account file is missing.
	ErrMissingGWSServiceAccountFile = fmt.Errorf("missing GWS service account file")
	// ErrMissingGWSUserEmail is returned when the GWS user email is missing.
	ErrMissingGWSUserEmail = fmt.Errorf("missing GWS user email")
)

// Config represents the configuration of the application.
type Config struct {
	ConfigFile string `mapstructure:"config-file"`
	IsLambda   bool
	Debug      bool

	LogLevel  string `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
	LogFormat string `mapstructure:"log_format" json:"log_format" yaml:"log_format"`

	GWSServiceAccountFile           string   `mapstructure:"gws_service_account_file" json:"gws_service_account_file" yaml:"gws_service_account_file"`
	GWSUserEmail                    string   `mapstructure:"gws_user_email" json:"gws_user_email" yaml:"gws_user_email"`
	GWSServiceAccountFileSecretName string   `mapstructure:"gws_service_account_file_secret_name" json:"gws_service_account_file_secret_name" yaml:"gws_service_account_file_secret_name"`
	GWSUserEmailSecretName          string   `mapstructure:"gws_user_email_secret_name" json:"gws_user_email_secret_name" yaml:"gws_user_email_secret_name"`
	GWSGroupsFilter                 []string `mapstructure:"gws_groups_filter" json:"gws_groups_filter" yaml:"gws_groups_filter"`
	GWSUsersFilter                  []string `mapstructure:"gws_users_filter" json:"gws_users_filter" yaml:"gws_users_filter"`
	GWSServiceAccountScopes         []string `mapstructure:"gws_service_account_scopes" json:"gws_service_account_scopes" yaml:"gws_service_account_scopes"`

	AWSSCIMEndpoint              string `mapstructure:"aws_scim_endpoint" json:"aws_scim_endpoint" yaml:"aws_scim_endpoint"`
	AWSSCIMAccessToken           string `mapstructure:"aws_scim_access_token" json:"aws_scim_access_token" yaml:"aws_scim_access_token"`
	AWSSCIMEndpointSecretName    string `mapstructure:"aws_scim_endpoint_secret_name" json:"aws_scim_endpoint_secret_name" yaml:"aws_scim_endpoint_secret_name"`
	AWSSCIMAccessTokenSecretName string `mapstructure:"aws_scim_access_token_secret_name" json:"aws_scim_access_token_secret_name" yaml:"aws_scim_access_token_secret_name"`

	AWSS3BucketName string `mapstructure:"aws_s3_bucket_name" json:"aws_s3_bucket_name" yaml:"aws_s3_bucket_name"`
	AWSS3BucketKey  string `mapstructure:"aws_s3_bucket_key" json:"aws_s3_bucket_key" yaml:"aws_s3_bucket_key"`

	// SyncMethod allow to defined the sync method used to get the user and groups from Google Workspace
	SyncMethod string `mapstructure:"sync_method" json:"sync_method" yaml:"sync_method"`

	// UseSecretsManager determines if we will use the AWS Secrets Manager secrets or program parameter values
	UseSecretsManager bool `mapstructure:"use_secrets_manager" json:"use_secrets_manager" yaml:"use_secrets_manager"`
}

// New returns a new Config
func New() Config {
	return Config{
		ConfigFile:                      DefaultConfigFile,
		IsLambda:                        DefaultIsLambda,
		Debug:                           DefaultDebug,
		LogLevel:                        DefaultLogLevel,
		LogFormat:                       DefaultLogFormat,
		GWSServiceAccountFile:           DefaultGWSServiceAccountFile,
		SyncMethod:                      DefaultSyncMethod,
		AWSS3BucketKey:                  DefaultAWSS3BucketKey,
		GWSServiceAccountFileSecretName: DefaultGWSServiceAccountFileSecretName,
		GWSUserEmailSecretName:          DefaultGWSUserEmailSecretName,
		AWSSCIMEndpointSecretName:       DefaultAWSSCIMEndpointSecretName,
		AWSSCIMAccessTokenSecretName:    DefaultAWSSCIMAccessTokenSecretName,
		UseSecretsManager:               DefaultUseSecretsManager,
		GWSServiceAccountScopes: []string{
			"https://www.googleapis.com/auth/admin.directory.group.readonly",
			"https://www.googleapis.com/auth/admin.directory.group.member.readonly",
			"https://www.googleapis.com/auth/admin.directory.user.readonly",
		},
	}
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true, "panic": true,
	}
	if !validLogLevels[c.LogLevel] {
		return ErrInvalidLogLevel
	}

	validLogFormats := map[string]bool{
		"text": true, "json": true,
	}
	if !validLogFormats[c.LogFormat] {
		return ErrInvalidLogFormat
	}

	if !c.UseSecretsManager {
		if c.AWSSCIMEndpoint == "" {
			return ErrMissingAWSSCIMEndpoint
		}
		if c.AWSSCIMAccessToken == "" {
			return ErrMissingAWSSCIMAccessToken
		}
		if c.GWSServiceAccountFile == "" {
			return ErrMissingGWSServiceAccountFile
		}
		if c.GWSUserEmail == "" {
			return ErrMissingGWSUserEmail
		}
	}

	return nil
}
