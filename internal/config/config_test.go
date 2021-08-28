package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := NewConfig()

	assert.NotNil(cfg)
	assert.NotNil(DefaultLogLevel)
	assert.NotNil(DefaultLogFormat)
	assert.NotNil(DefaultDebug)
	assert.NotNil(DefaultServiceAccountFile)

	assert.Equal(cfg.LogLevel, DefaultLogLevel)
	assert.Equal(cfg.LogFormat, DefaultLogFormat)
	assert.Equal(cfg.Debug, DefaultDebug)
	assert.Equal(cfg.ServiceAccountFile, DefaultServiceAccountFile)
}

func TestConfig_toJSON(t *testing.T) {
	type fields struct {
		Debug              bool
		LogLevel           string
		LogFormat          string
		ServiceAccountFile string
		UserEmail          string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name:   "default",
			fields: fields{Debug: false, LogLevel: "info", LogFormat: "text", ServiceAccountFile: "default"},
			want: []byte(`{
  "IsLambda": false,
  "Debug": false,
  "LogLevel": "info",
  "LogFormat": "text",
  "ServiceAccountFile": "default",
  "UserEmail": "",
  "SCIMEndpoint": "",
  "SCIMAccessToken": ""
}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Debug:              tt.fields.Debug,
				LogLevel:           tt.fields.LogLevel,
				LogFormat:          tt.fields.LogFormat,
				ServiceAccountFile: tt.fields.ServiceAccountFile,
				UserEmail:          tt.fields.UserEmail,
			}
			if got := c.toJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.toJSON() =\n%s, \n want %s", got, tt.want)
			}
		})
	}
}
