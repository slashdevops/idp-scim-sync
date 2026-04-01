package config

import (
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
	assert.Equal(cfg.AWSSCIMEndpointSecretName, DefaultAWSSCIMEndpointSecretName)
	assert.Equal(cfg.AWSSCIMAccessTokenSecretName, DefaultAWSSCIMAccessTokenSecretName)
	assert.Equal(cfg.UseSecretsManager, DefaultUseSecretsManager)
}

func validConfig() Config {
	return Config{
		LogLevel:              "info",
		LogFormat:             "json",
		AWSSCIMEndpoint:       "https://scim.example.com",
		AWSSCIMAccessToken:    "token",
		GWSServiceAccountFile: "credentials.json",
		GWSUserEmail:          "admin@example.com",
	}
}

func TestValidate(t *testing.T) {
	t.Run("valid config passes", func(t *testing.T) {
		cfg := validConfig()
		assert.NoError(t, cfg.Validate())
	})

	t.Run("invalid log level", func(t *testing.T) {
		cfg := validConfig()
		cfg.LogLevel = "invalid"
		err := cfg.Validate()
		assert.ErrorIs(t, err, ErrInvalidLogLevel)
	})

	t.Run("all valid log levels", func(t *testing.T) {
		for _, level := range []string{"debug", "info", "warn", "error", "fatal", "panic"} {
			cfg := validConfig()
			cfg.LogLevel = level
			assert.NoError(t, cfg.Validate(), "expected %q to be valid", level)
		}
	})

	t.Run("invalid log format", func(t *testing.T) {
		cfg := validConfig()
		cfg.LogFormat = "yaml"
		err := cfg.Validate()
		assert.ErrorIs(t, err, ErrInvalidLogFormat)
	})

	t.Run("valid log formats", func(t *testing.T) {
		for _, format := range []string{"text", "json"} {
			cfg := validConfig()
			cfg.LogFormat = format
			assert.NoError(t, cfg.Validate(), "expected %q to be valid", format)
		}
	})

	t.Run("missing AWS SCIM endpoint", func(t *testing.T) {
		cfg := validConfig()
		cfg.AWSSCIMEndpoint = ""
		err := cfg.Validate()
		assert.ErrorIs(t, err, ErrMissingAWSSCIMEndpoint)
	})

	t.Run("missing AWS SCIM access token", func(t *testing.T) {
		cfg := validConfig()
		cfg.AWSSCIMAccessToken = ""
		err := cfg.Validate()
		assert.ErrorIs(t, err, ErrMissingAWSSCIMAccessToken)
	})

	t.Run("missing GWS service account file", func(t *testing.T) {
		cfg := validConfig()
		cfg.GWSServiceAccountFile = ""
		err := cfg.Validate()
		assert.ErrorIs(t, err, ErrMissingGWSServiceAccountFile)
	})

	t.Run("missing GWS user email", func(t *testing.T) {
		cfg := validConfig()
		cfg.GWSUserEmail = ""
		err := cfg.Validate()
		assert.ErrorIs(t, err, ErrMissingGWSUserEmail)
	})

	t.Run("secrets manager skips credential validation", func(t *testing.T) {
		cfg := validConfig()
		cfg.UseSecretsManager = true
		cfg.AWSSCIMEndpoint = ""
		cfg.AWSSCIMAccessToken = ""
		cfg.GWSServiceAccountFile = ""
		cfg.GWSUserEmail = ""
		assert.NoError(t, cfg.Validate())
	})

	t.Run("valid sync_user_fields", func(t *testing.T) {
		cfg := validConfig()
		cfg.SyncUserFields = []string{"phoneNumbers", "addresses", "enterpriseData"}
		assert.NoError(t, cfg.Validate())
	})

	t.Run("invalid sync_user_fields", func(t *testing.T) {
		cfg := validConfig()
		cfg.SyncUserFields = []string{"phoneNumbers", "invalidField"}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalidField")
	})

	t.Run("empty string in sync_user_fields is ignored", func(t *testing.T) {
		cfg := validConfig()
		cfg.SyncUserFields = []string{"", "phoneNumbers", ""}
		assert.NoError(t, cfg.Validate())
	})

	t.Run("only empty strings in sync_user_fields passes", func(t *testing.T) {
		cfg := validConfig()
		cfg.SyncUserFields = []string{""}
		assert.NoError(t, cfg.Validate())
	})
}
