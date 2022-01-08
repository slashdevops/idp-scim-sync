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
	assert.Equal(cfg.GWSServiceAccountFileSecretARN, DefaultGWSServiceAccountFileSecretARN)
	assert.Equal(cfg.GWSUserEmailSecretARN, DefaultGWSUserEmailSecretARN)
	assert.Equal(cfg.AWSSCIMEndpointSecretARN, DefaultAWSSCIMEndpointSecretARN)
	assert.Equal(cfg.AWSSCIMAccessTokenSecretARN, DefaultAWSSCIMAccessTokenSecretARN)
}
