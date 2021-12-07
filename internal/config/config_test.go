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
  "AWSS3BucketKey": "",
  "SyncMethod": "groups",
  "StateEnabled": false
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
loglevel: info
logformat: text
gwsserviceaccountfile: credentials.json
gwsuseremail: ""
gwsserviceaccountfilesecretname: IDPSCIM_GWSServiceAccountFile
gwsuseremailsecretname: IDPSCIM_GWSUserEmail
gwsgroupsfilter: []
gwsusersfilter: []
scimendpoint: ""
scimaccesstoken: ""
scimendpointsecretname: IDPSCIM_SCIMEndpoint
scimaccesstokensecretname: IDPSCIM_SCIMAccessToken
awss3bucketname: ""
awss3bucketkey: ""
syncmethod: groups
stateenabled: false
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
