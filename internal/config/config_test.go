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
	assert.NotNil(DefaultLogLevel)
	assert.NotNil(DefaultLogFormat)
	assert.NotNil(DefaultDebug)
	assert.NotNil(DefaultServiceAccountFile)
	assert.NotNil(DefaultSyncMethod)

	assert.Equal(cfg.IsLambda, false)
	assert.Equal(cfg.Debug, DefaultDebug)
	assert.Equal(cfg.LogLevel, DefaultLogLevel)
	assert.Equal(cfg.LogFormat, DefaultLogFormat)
	assert.Equal(cfg.ServiceAccountFile, DefaultServiceAccountFile)
	assert.Equal(cfg.SyncMethod, "groups")
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
  "ServiceAccountFile": "credentials.json",
  "UserEmail": "",
  "SCIMEndpoint": "",
  "SCIMAccessToken": "",
  "SyncMethod": "groups"
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
serviceaccountfile: credentials.json
useremail: ""
scimendpoint: ""
scimaccesstoken: ""
syncmethod: groups
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
