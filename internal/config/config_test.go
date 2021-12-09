package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	cfg := New()

	assert.NotNil(cfg)

	assert.False(cfg.IsLambda)
	assert.Equal(cfg.Debug, DefaultDebug)
	assert.Equal(cfg.LogLevel, DefaultLogLevel)
	assert.Equal(cfg.LogFormat, DefaultLogFormat)
	assert.Equal(cfg.GWSServiceAccountFile, DefaultGWSServiceAccountFile)
	assert.Equal(cfg.SyncMethod, DefaultSyncMethod)
	assert.Equal(cfg.GWSServiceAccountFileSecretName, DefaultGWSServiceAccountFileSecretName)
	assert.Equal(cfg.GWSUserEmailSecretName, DefaultGWSUserEmailSecretName)
	assert.Equal(cfg.SCIMEndpointSecretName, DefaultSCIMEndpointSecretName)
	assert.Equal(cfg.SCIMAccessTokenSecretName, DefaultSCIMAccessTokenSecretName)
}

func TestConfig_toJSON(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   []byte
	}{
		{
			name:   "default values",
			config: New(),
			want: []byte(`{
  "IsLambda": false,
  "Debug": false,
  "LogLevel": "info",
  "LogFormat": "text",
  "GWSServiceAccountFile": "credentials.json",
  "GWSUserEmail": "",
  "GWSServiceAccountFileSecretName": "IDPSCIM_GWSServiceAccountFile",
  "GWSUserEmailSecretName": "IDPSCIM_GWSUserEmail",
  "GWSGroupsFilter": null,
  "GWSUsersFilter": null,
  "SCIMEndpoint": "",
  "SCIMAccessToken": "",
  "SCIMEndpointSecretName": "IDPSCIM_SCIMEndpoint",
  "SCIMAccessTokenSecretName": "IDPSCIM_SCIMAccessToken",
  "AWSS3BucketName": "",
  "AWSS3BucketKey": "state.json",
  "SyncMethod": "groups",
  "DisableState": false
}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.toJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.toJSON() =\n%s, \n want %s", got, tt.want)
			}
		})
	}
}

func TestConfig_toYAML(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   []byte
	}{
		{
			name:   "default values",
			config: New(),
			want: []byte(`islambda: false
debug: false
log_level: info
log_format: text
gws_service_account_file: credentials.json
gws_user_email: ""
gws_service_account_file_secret_name: IDPSCIM_GWSServiceAccountFile
gws_user_email_secret_name: IDPSCIM_GWSUserEmail
gws_groups_filter: []
gws_users_filter: []
scim_endpoint: ""
scim_access_token: ""
scim_endpoint_secret_name: IDPSCIM_SCIMEndpoint
scim_access_token_secret_name: IDPSCIM_SCIMAccessToken
aws_s3_bucket_name: ""
aws_s3_bucket_key: state.json
sync_method: groups
disable_state: false
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.toYAML(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.toYAML() = %s, want %s", got, tt.want)
			}
		})
	}
}
