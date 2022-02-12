package aws

import (
	"context"
	"os"
	"testing"
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
