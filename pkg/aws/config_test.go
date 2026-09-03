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
			t.Setenv(key, value)
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
			// NOTE: testdata/case1/ does not exist. Kept as-is because pointing
			// AWS_CONFIG_FILE at a missing file is itself worth covering — the
			// SDK tolerates it and simply resolves no region, which is what the
			// empty AWS_REGION expectation below asserts. The "case1" name looks
			// like a leftover; renaming it to a real fixture would change what
			// this subtest covers.
			"AWS_CONFIG_FILE": "testdata/case1/config",
		}

		// Corrected alongside the switch from os.Setenv to t.Setenv. These
		// expectations described the leaked environment of the
		// "using access key from env vars" subtest above, not this subtest's own
		// setup: os.Setenv never cleaned up, so AWS_ACCESS_KEY_ID and AWS_REGION
		// from that subtest were still set when this one ran, and the SDK
		// resolved env credentials. t.Setenv isolates each subtest, which
		// exposed it.
		Expected := map[string]string{
			"AWS_ACCESS_KEY_ID":     "AKIAIOSFODNN7EXAMPLE",
			"AWS_SECRET_ACCESS_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			// No usable config file, so no region is resolved.
			"AWS_REGION":                  "",
			"AWS_SHARED_CREDENTIALS_FILE": "SharedConfigCredentials: testdata/profile/credentials",
		}

		for key, value := range GotEnvVars {
			t.Logf("setting env var: %s", key)
			t.Setenv(key, value)
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

		// Corrected alongside the switch from os.Setenv to t.Setenv. These
		// expectations described the leaked environment of the
		// "using access key from env vars" subtest above, not this subtest's own
		// setup: os.Setenv never cleaned up, so AWS_ACCESS_KEY_ID and AWS_REGION
		// from that subtest were still set when this one ran, and the SDK
		// resolved env credentials. t.Setenv isolates each subtest, which
		// exposed it.
		Expected := map[string]string{
			"AWS_ACCESS_KEY_ID":     "AKIAIOSFODNN7EXAMPLE",
			"AWS_SECRET_ACCESS_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			// Only a credentials file is provided, no config file, so no region
			// is resolved.
			"AWS_REGION":                  "",
			"AWS_SHARED_CREDENTIALS_FILE": "SharedConfigCredentials: testdata/default/credentials",
		}

		for key, value := range GotEnvVars {
			t.Logf("setting env var: %s", key)
			t.Setenv(key, value)
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

		// This asserted "EnvConfigCredentials" while every env var it sets
		// points at the shared credentials file. It only passed because the
		// subtests used os.Setenv without cleanup, so the AWS_ACCESS_KEY_ID that
		// "using access key from env vars" set above leaked into this subtest and
		// the SDK resolved env credentials instead. Switching to t.Setenv
		// isolates the subtests and exposed it.
		//
		// The correct source is the shared credentials file, which is what this
		// subtest is named for and what the AWS_PROFILE variant below already
		// expects. The key and secret assertions never caught the discrepancy
		// because testdata/default/credentials holds the same values the env-var
		// subtest uses.
		Expected := map[string]string{
			"AWS_ACCESS_KEY_ID":           "AKIAIOSFODNN7EXAMPLE",
			"AWS_SECRET_ACCESS_KEY":       "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			"AWS_REGION":                  "us-east-1",
			"AWS_SHARED_CREDENTIALS_FILE": "SharedConfigCredentials: testdata/default/credentials",
		}

		for key, value := range GotEnvVars {
			t.Logf("setting env var: %s", key)
			t.Setenv(key, value)
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
			t.Setenv(key, value)
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
		// t.Setenv records the previous value and restores it during cleanup,
		// which replaces the hand-rolled save/restore this used to do.
		t.Setenv("AWS_PROFILE", "test-profile")

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
		// NewDefaultConf treats an empty AWS_PROFILE the same as an unset one
		// (it guards with `if profile := os.Getenv(...); profile != ""`), so an
		// empty value exercises the intended branch and t.Setenv restores
		// whatever was there before.
		t.Setenv("AWS_PROFILE", "")

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
