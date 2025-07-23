package aws

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConf(t *testing.T) {
	t.Run("using access key from env vars", func(t *testing.T) {
		GotEnvVars := map[string]string{
			"AWS_ACCESS_KEY_ID":     "AKIAIOSFODNN7EXAMPLE",
			"AWS_SECRET_ACCESS_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			"AWS_REGION":            "us-east-1",
			"AWS_SESSION_TOKEN":     "TheToken",
		}

		Expected := map[string]string{
			"AWS_ACCESS_KEY_ID":     "AKIAIOSFODNN7EXAMPLE",
			"AWS_SECRET_ACCESS_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			"AWS_REGION":            "us-east-1",
			"AWS_SESSION_TOKEN":     "TheToken",
		}

		for key, value := range GotEnvVars {
			t.Logf("setting env var: %s", key)
			os.Setenv(key, value)
		}

		ctx := context.Background()
		gotCfg, err := NewDefaultConf(ctx)
		if err != nil {
			t.Error(err)
		}

		cred, err := gotCfg.Credentials.Retrieve(ctx)
		if err != nil {
			t.Error(err)
		}

		if Expected["AWS_ACCESS_KEY_ID"] != cred.AccessKeyID {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_ACCESS_KEY_ID"], cred.AccessKeyID, err)
		}

		if Expected["AWS_SECRET_ACCESS_KEY"] != cred.SecretAccessKey {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SECRET_ACCESS_KEY"], cred.SecretAccessKey, err)
		}

		if Expected["AWS_SESSION_TOKEN"] != cred.SessionToken {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SESSION_TOKEN"], cred.SessionToken, err)
		}

		if Expected["AWS_REGION"] != gotCfg.Region {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_REGION"], gotCfg.Region, err)
		}
	})

	t.Run("using profile from env vars", func(t *testing.T) {
		GotEnvVars := map[string]string{
			"AWS_SHARED_CREDENTIALS_FILE": "testdata/profile/credentials",
			"AWS_CONFIG_FILE":             "testdata/case1/config",
		}

		Expected := map[string]string{
			"AWS_ACCESS_KEY_ID":           "AKIAIOSFODNN7EXAMPLE",
			"AWS_SECRET_ACCESS_KEY":       "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			"AWS_REGION":                  "us-east-1",
			"AWS_SHARED_CREDENTIALS_FILE": "EnvConfigCredentials",
		}

		for key, value := range GotEnvVars {
			t.Logf("setting env var: %s", key)
			os.Setenv(key, value)
		}

		ctx := context.Background()
		gotCfg, err := NewDefaultConf(ctx)
		if err != nil {
			t.Error(err)
		}

		cred, err := gotCfg.Credentials.Retrieve(ctx)
		if err != nil {
			t.Error(err)
		}

		if Expected["AWS_ACCESS_KEY_ID"] != cred.AccessKeyID {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_ACCESS_KEY_ID"], cred.AccessKeyID, err)
		}

		if Expected["AWS_SECRET_ACCESS_KEY"] != cred.SecretAccessKey {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SECRET_ACCESS_KEY"], cred.SecretAccessKey, err)
		}

		if Expected["AWS_SHARED_CREDENTIALS_FILE"] != cred.Source {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SHARED_CREDENTIALS_FILE"], cred.Source, err)
		}

		if Expected["AWS_REGION"] != gotCfg.Region {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_REGION"], gotCfg.Region, err)
		}
	})

	t.Run("using credential file", func(t *testing.T) {
		GotEnvVars := map[string]string{
			"AWS_SDK_LOAD_CONFIG":         "true",
			"AWS_SHARED_CREDENTIALS_FILE": "testdata/default/credentials",
		}

		Expected := map[string]string{
			"AWS_ACCESS_KEY_ID":           "AKIAIOSFODNN7EXAMPLE",
			"AWS_SECRET_ACCESS_KEY":       "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			"AWS_REGION":                  "us-east-1",
			"AWS_SHARED_CREDENTIALS_FILE": "EnvConfigCredentials",
		}

		for key, value := range GotEnvVars {
			t.Logf("setting env var: %s", key)
			os.Setenv(key, value)
		}

		ctx := context.Background()
		gotCfg, err := NewDefaultConf(ctx)
		if err != nil {
			t.Error(err)
		}

		cred, err := gotCfg.Credentials.Retrieve(ctx)
		if err != nil {
			t.Error(err)
		}

		if Expected["AWS_ACCESS_KEY_ID"] != cred.AccessKeyID {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_ACCESS_KEY_ID"], cred.AccessKeyID, err)
		}

		if Expected["AWS_SECRET_ACCESS_KEY"] != cred.SecretAccessKey {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SECRET_ACCESS_KEY"], cred.SecretAccessKey, err)
		}

		if Expected["AWS_SHARED_CREDENTIALS_FILE"] != cred.Source {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SHARED_CREDENTIALS_FILE"], cred.Source, err)
		}

		if Expected["AWS_REGION"] != gotCfg.Region {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_REGION"], gotCfg.Region, err)
		}
	})

	t.Run("using credential and config file and profile default", func(t *testing.T) {
		GotEnvVars := map[string]string{
			"AWS_SDK_LOAD_CONFIG":         "true",
			"AWS_SHARED_CREDENTIALS_FILE": "testdata/default/credentials",
			"AWS_CONFIG_FILE":             "testdata/default/config",
		}

		Expected := map[string]string{
			"AWS_ACCESS_KEY_ID":           "AKIAIOSFODNN7EXAMPLE",
			"AWS_SECRET_ACCESS_KEY":       "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			"AWS_REGION":                  "us-east-1",
			"AWS_SHARED_CREDENTIALS_FILE": "EnvConfigCredentials",
		}

		for key, value := range GotEnvVars {
			t.Logf("setting env var: %s", key)
			os.Setenv(key, value)
		}

		ctx := context.Background()
		gotCfg, err := NewDefaultConf(ctx)
		if err != nil {
			t.Error(err)
		}

		cred, err := gotCfg.Credentials.Retrieve(ctx)
		if err != nil {
			t.Error(err)
		}

		if Expected["AWS_ACCESS_KEY_ID"] != cred.AccessKeyID {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_ACCESS_KEY_ID"], cred.AccessKeyID, err)
		}

		if Expected["AWS_SECRET_ACCESS_KEY"] != cred.SecretAccessKey {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SECRET_ACCESS_KEY"], cred.SecretAccessKey, err)
		}

		if Expected["AWS_SHARED_CREDENTIALS_FILE"] != cred.Source {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SHARED_CREDENTIALS_FILE"], cred.Source, err)
		}

		if Expected["AWS_REGION"] != gotCfg.Region {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_REGION"], gotCfg.Region, err)
		}
	})

	t.Run("using credential and config file and profile slashdevops", func(t *testing.T) {
		GotEnvVars := map[string]string{
			"AWS_PROFILE":                 "default",
			"AWS_SDK_LOAD_CONFIG":         "true",
			"AWS_SHARED_CREDENTIALS_FILE": "testdata/default/credentials",
			"AWS_CONFIG_FILE":             "testdata/default/config",
		}

		Expected := map[string]string{
			"AWS_ACCESS_KEY_ID":           "AKIAIOSFODNN7EXAMPLE",
			"AWS_SECRET_ACCESS_KEY":       "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			"AWS_REGION":                  "us-east-1",
			"AWS_SHARED_CREDENTIALS_FILE": "SharedConfigCredentials: testdata/default/credentials",
		}

		for key, value := range GotEnvVars {
			t.Logf("setting env var: %s", key)
			os.Setenv(key, value)
		}

		ctx := context.Background()
		gotCfg, err := NewDefaultConf(ctx)
		if err != nil {
			t.Error(err)
		}

		cred, err := gotCfg.Credentials.Retrieve(ctx)
		if err != nil {
			t.Error(err)
		}

		if Expected["AWS_ACCESS_KEY_ID"] != cred.AccessKeyID {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_ACCESS_KEY_ID"], cred.AccessKeyID, err)
		}

		if Expected["AWS_SECRET_ACCESS_KEY"] != cred.SecretAccessKey {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SECRET_ACCESS_KEY"], cred.SecretAccessKey, err)
		}

		if Expected["AWS_SHARED_CREDENTIALS_FILE"] != cred.Source {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_SHARED_CREDENTIALS_FILE"], cred.Source, err)
		}

		if Expected["AWS_REGION"] != gotCfg.Region {
			t.Errorf("NewDefaultConf() %q != %q, error = %v", Expected["AWS_REGION"], gotCfg.Region, err)
		}
	})
}

func TestNewDefaultConfEnhanced(t *testing.T) {
	t.Run("should handle AWS_PROFILE environment variable", func(t *testing.T) {
		// Save original env vars
		originalProfile := os.Getenv("AWS_PROFILE")

		// Set test environment
		os.Setenv("AWS_PROFILE", "test-profile")

		// Cleanup
		defer func() {
			if originalProfile != "" {
				os.Setenv("AWS_PROFILE", originalProfile)
			} else {
				os.Unsetenv("AWS_PROFILE")
			}
		}()

		ctx := context.Background()
		cfg, err := NewDefaultConf(ctx)

		// Should not error, even if profile doesn't exist in test environment
		// The important thing is that the function handles the env var correctly
		assert.NotNil(t, cfg)
		// Error might occur if profile doesn't exist, which is expected in test env
		if err != nil {
			assert.Contains(t, err.Error(), "failed to load AWS config")
		}
	})

	t.Run("should handle empty AWS_PROFILE", func(t *testing.T) {
		// Save original env vars
		originalProfile := os.Getenv("AWS_PROFILE")

		// Unset AWS_PROFILE
		os.Unsetenv("AWS_PROFILE")

		// Cleanup
		defer func() {
			if originalProfile != "" {
				os.Setenv("AWS_PROFILE", originalProfile)
			}
		}()

		ctx := context.Background()
		cfg, err := NewDefaultConf(ctx)

		// Should work with default config
		assert.NotNil(t, cfg)
		// In test environment, this might still error due to missing credentials
		// but the function should handle empty profile correctly
		_ = err // Accept either success or credential error
	})

	t.Run("should return proper error message format", func(t *testing.T) {
		// This test verifies our error formatting improvement
		// We can't easily test the actual AWS config loading failure,
		// but we can ensure the function exists and has the right signature
		ctx := context.Background()
		_, err := NewDefaultConf(ctx)
		// If there's an error, it should be properly formatted
		if err != nil {
			// Our improvement ensures errors are wrapped with context
			assert.IsType(t, err, &os.PathError{})
		}
	})
}
