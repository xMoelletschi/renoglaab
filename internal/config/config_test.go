//nolint:lll
package config

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const defaultBranchRegex = `^renovate/automerge$`

func TestNewConfig(t *testing.T) {
	t.Setenv("CONFIG_PATH", "config.js")
	t.Setenv("LABELS", "renovate")
	t.Setenv("APPROVE_COMMENT", "Approving merge request! :ship:")
	t.Setenv("GITLAB_API_TOKEN", "glpat-dafsgretegsfaf")
	t.Setenv("GITLAB_URL", "https://test.gitlab.com")
	t.Setenv("LOG_LEVEL", "warn")

	config := NewConfig()

	assert.Equal(t, "config.js", config.ConfigPath, "Config does not match")
	assert.Equal(t, []string{"renovate"}, config.Labels, "Labels does not match")
	assert.Equal(t, "Approving merge request! :ship:", config.Comment, "ApproveComment does not match'")
	assert.Equal(t, "glpat-dafsgretegsfaf", config.GitLabAPIToken, "GitLabAPIToken does not match")
	assert.Equal(t, "https://test.gitlab.com", config.GitLabURL, "GitLabURL does not match")
	assert.Equal(t, regexp.MustCompile(defaultBranchRegex).String(), config.AllowedBranchRegex.String(), "AllowedBranchRegex does not match")
	assert.Equal(t, logrus.WarnLevel, config.LogLevel)
}

func TestNewConfigWithInvalidLogLevel(t *testing.T) {
	t.Setenv("LOG_LEVEL", "invalid")

	config := NewConfig()

	assert.Equal(t, logrus.InfoLevel, config.LogLevel, "LogLevel should default to InfoLevel when an invalid log level is provided")
}

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	defaultValue := "default_value"

	value := getEnv(key, defaultValue)
	assert.Equal(t, defaultValue, value, "Expected default value when environment variable is not set")

	expectedValue := "set_value"
	t.Setenv(key, expectedValue)
	value = getEnv(key, defaultValue)
	assert.Equal(t, expectedValue, value, "Expected environment variable value when it is set")
}

func TestGetEnvAsSlice(t *testing.T) {
	key := "TEST_ENV_SLICE"
	defaultValue := "default1,default2"

	value := getEnvAsSlice(key, defaultValue)
	expected := strings.Split(defaultValue, ",")
	assert.Equal(t, expected, value, "Expected default value slice when environment variable is not set")

	expectedValue := "value1,value2,value3"
	t.Setenv(key, expectedValue)
	value = getEnvAsSlice(key, defaultValue)
	expected = strings.Split(expectedValue, ",")
	assert.Equal(t, expected, value, "Expected environment variable value slice when it is set")
}

func TestGetEnvAsBool(t *testing.T) {
	key := "TEST_ENV_BOOL"
	defaultValue := false

	value := getEnvAsBool(key, defaultValue)
	assert.Equal(t, defaultValue, value, "Expected default value when environment variable is not set")

	t.Setenv(key, "true")
	value = getEnvAsBool(key, defaultValue)
	assert.True(t, value, "Expected true when environment variable is set to 'true'")

	t.Setenv(key, "false")
	value = getEnvAsBool(key, defaultValue)
	assert.False(t, value, "Expected false when environment variable is set to 'false'")

	t.Setenv(key, "some_random_value")
	value = getEnvAsBool(key, defaultValue)
	assert.False(t, value, "Expected false when environment variable is set to an arbitrary string")
}
func TestPrintConfig(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("CONFIG_PATH", "config.js")
	t.Setenv("GITLAB_URL", "https://test.gitlab.com")
	t.Setenv("LABELS", "renovate")
	t.Setenv("FILTER_BY_BRANCH", "true")
	t.Setenv("FILTER_BY_SUCCEEDED_PIPELINE", "true")
	t.Setenv("FILTER_BY_PIPELINE_WITHOUT_WARNINGS", "true")
	t.Setenv("FILTER_BY_LABELS", "true")
	t.Setenv("FILTER_BY_AUTHOR_USERNAME", "true")
	t.Setenv("AUTHOR_USERNAME", "renovate-bot")
	t.Setenv("ADD_COMMENT", "true")
	t.Setenv("APPROVE_COMMENT", "Approving merge request! :ship:")

	config := NewConfig()

	// Capture the log output
	var logOutput strings.Builder

	logrus.SetOutput(&logOutput)
	defer logrus.SetOutput(os.Stderr)

	config.PrintConfig()

	logContent := logOutput.String()
	assert.Contains(t, logContent, "LogLevel=debug", "Expected LogLevel in log output")
	assert.Contains(t, logContent, "ConfigPath=config.js", "Expected ConfigPath in log output")
	assert.Contains(t, logContent, "GitLabURL=\"https://test.gitlab.com\"", "Expected GitLabURL in log output")
	assert.Contains(t, logContent, "Labels=\"[renovate]\"", "Expected Labels in log output")
	assert.Contains(t, logContent, "FilterByBranch=true", "Expected FilterByBranch in log output")
	assert.Contains(t, logContent, "FilterBySucceededPipeline=true", "Expected FilterBySucceededPipeline in log output")
	assert.Contains(t, logContent, "FilterByPipelineWithoutWarnings=true", "Expected FilterByPipelineWithoutWarnings in log output")
	assert.Contains(t, logContent, "FilterByLabels=true", "Expected FilterByLabels in log output")
	assert.Contains(t, logContent, "FilterByAuthorUsername=true", "Expected FilterByAuthorUsername in log output")
	assert.Contains(t, logContent, "AuthorUsername=renovate-bot", "Expected AuthorUsername in log output")
	assert.Contains(t, logContent, "AddComment=true", "Expected AddComment in log output")
}
