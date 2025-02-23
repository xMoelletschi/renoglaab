//nolint:lll
package config

import (
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	ConfigPath                      string
	LogLevel                        logrus.Level
	GitLabAPIToken                  string
	GitLabURL                       string
	FilterByAuthorUsername          bool
	AuthorUsername                  string
	FilterByLabels                  bool
	Labels                          []string
	FilterByBranch                  bool
	AllowedBranchRegex              string
	AllowedBranchRegexCompiled      *regexp.Regexp
	FilterBySucceededPipeline       bool
	FilterByPipelineWithoutWarnings bool
	AddComment                      bool
	Comment                         string
	Approve                         string
}

// getDefaultConfig returns the default configuration values.
func getDefaultConfig() Config {
	return Config{
		ConfigPath:                      "$CI_PROJECT_DIR/config.js",
		LogLevel:                        logrus.InfoLevel,
		GitLabURL:                       "https://gitlab.com",
		FilterByAuthorUsername:          true,
		AuthorUsername:                  "renovate-bot",
		FilterByLabels:                  true,
		Labels:                          []string{"renovate"},
		FilterByBranch:                  true,
		AllowedBranchRegex:              `renovate/automerge`,
		FilterBySucceededPipeline:       true,
		FilterByPipelineWithoutWarnings: true,
		AddComment:                      true,
		Comment:                         "Approving merge request! :ship:",
		Approve:                         "/approve",
	}
}

func NewConfig() *Config {
	cfg := getDefaultConfig() // Use default values

	cfg.ConfigPath = os.ExpandEnv(getEnv("CONFIG_PATH", cfg.ConfigPath))
	cfg.LogLevel = mustParseLogLevel(getEnv("LOG_LEVEL", cfg.LogLevel.String()))
	cfg.GitLabAPIToken = getEnv("GITLAB_API_TOKEN", "")
	cfg.GitLabURL = getEnv("GITLAB_URL", cfg.GitLabURL)
	cfg.FilterByAuthorUsername = getEnvAsBool("FILTER_BY_AUTHOR_USERNAME", cfg.FilterByAuthorUsername)
	cfg.AuthorUsername = getEnv("AUTHOR_USERNAME", cfg.AuthorUsername)
	cfg.FilterByLabels = getEnvAsBool("FILTER_BY_LABELS", cfg.FilterByLabels)
	cfg.Labels = getEnvAsSlice("LABELS", strings.Join(cfg.Labels, ","))
	cfg.FilterByBranch = getEnvAsBool("FILTER_BY_BRANCH", cfg.FilterByBranch)
	cfg.AllowedBranchRegex = getEnv("ALLOWED_BRANCH_REGEX", cfg.AllowedBranchRegex)
	cfg.AllowedBranchRegexCompiled = mustCompileRegex(cfg.AllowedBranchRegex)
	cfg.FilterBySucceededPipeline = getEnvAsBool("FILTER_BY_SUCCEEDED_PIPELINE", cfg.FilterBySucceededPipeline)
	cfg.FilterByPipelineWithoutWarnings = getEnvAsBool("FILTER_BY_PIPELINE_WITHOUT_WARNINGS", cfg.FilterByPipelineWithoutWarnings)
	cfg.AddComment = getEnvAsBool("ADD_COMMENT", cfg.AddComment)
	cfg.Comment = getEnv("COMMENT", cfg.Comment)
	cfg.Approve = getEnv("APPROVE", cfg.Approve)

	configureLogging(&cfg)

	cfg.PrintConfig()

	return &cfg
}

// PrintConfig logs the configuration when debug is enabled.
func (c *Config) PrintConfig() {
	if logrus.GetLevel() == logrus.DebugLevel {
		logrus.WithFields(logrus.Fields{
			"ConfigPath":                      c.ConfigPath,
			"LogLevel":                        c.LogLevel.String(),
			"GitLabURL":                       c.GitLabURL,
			"FilterByAuthorUsername":          c.FilterByAuthorUsername,
			"AuthorUsername":                  c.AuthorUsername,
			"FilterByLabels":                  c.FilterByLabels,
			"Labels":                          c.Labels,
			"FilterByBranch":                  c.FilterByBranch,
			"AllowedBranchRegex":              c.AllowedBranchRegex,
			"AllowedBranchRegexCompiled":      c.AllowedBranchRegexCompiled,
			"FilterBySucceededPipeline":       c.FilterBySucceededPipeline,
			"FilterByPipelineWithoutWarnings": c.FilterByPipelineWithoutWarnings,
			"AddComment":                      c.AddComment,
			"Comment":                         c.Comment,
		}).Debug("Loaded Configuration")
	}
}

func configureLogging(cfg *Config) {
	logFormat := getEnv("LOG_FORMAT", "text")
	if logFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	logrus.SetLevel(cfg.LogLevel)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		switch strings.ToLower(value) {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}

	return defaultValue
}

func getEnvAsSlice(key, defaultValue string) []string {
	valueStr := getEnv(key, defaultValue)
	if valueStr == "" {
		return nil
	}

	return strings.Split(valueStr, ",")
}

func mustCompileRegex(pattern string) *regexp.Regexp {
	if !strings.HasPrefix(pattern, "^") {
		pattern = "^" + pattern
	}

	if !strings.HasSuffix(pattern, "$") {
		pattern += "$"
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		logrus.Fatalf("Invalid regex pattern: %v", err)
	}

	return re
}

func mustParseLogLevel(levelStr string) logrus.Level {
	logLevel, err := logrus.ParseLevel(levelStr)
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	return logLevel
}
