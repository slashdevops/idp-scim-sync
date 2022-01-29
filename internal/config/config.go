package config

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

	AWSSCIMEndpoint              string `mapstructure:"aws_scim_endpoint" json:"aws_scim_endpoint" yaml:"aws_scim_endpoint"`
	AWSSCIMAccessToken           string `mapstructure:"aws_scim_access_token" json:"aws_scim_access_token" yaml:"aws_scim_access_token"`
	AWSSCIMEndpointSecretName    string `mapstructure:"aws_scim_endpoint_secret_name" json:"aws_scim_endpoint_secret_name" yaml:"aws_scim_endpoint_secret_name"`
	AWSSCIMAccessTokenSecretName string `mapstructure:"aws_scim_access_token_secret_name" json:"aws_scim_access_token_secret_name" yaml:"aws_scim_access_token_secret_name"`

	AWSS3BucketName string `mapstructure:"aws_s3_bucket_name" json:"aws_s3_bucket_name" yaml:"aws_s3_bucket_name"`
	AWSS3BucketKey  string `mapstructure:"aws_s3_bucket_key" json:"aws_s3_bucket_key" yaml:"aws_s3_bucket_key"`

	// SyncMethod allow to defined the sync method used to get the user and groups from Google Workspace
	SyncMethod string `mapstructure:"sync_method" json:"sync_method" yaml:"sync_method"`
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
	}
}
