//nolint:funlen
package gitlab_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gl "github.com/xMoelletschi/renoglaab/internal/gitlab"
)

func TestExtractFromEnv(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expected    []string
		expectError bool
	}{
		{
			name:        "Valid repositories",
			envValue:    "repo1 repo2 repo3",
			expected:    []string{"repo1", "repo2", "repo3"},
			expectError: false,
		},
		{
			name:        "Repositories with flags",
			envValue:    "repo1 --autodiscover=true repo2 --autodiscover-filter=repo3",
			expected:    []string{"repo1", "repo2"},
			expectError: false,
		},
		{
			name:        "Only flags",
			envValue:    "--autodiscover=true --autodiscover-filter=repo3",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Empty environment variable",
			envValue:    "",
			expected:    nil,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("RENOVATE_EXTRA_FLAGS", tt.envValue)

			repositories, err := gl.ExtractFromEnv()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, repositories)
			}
		})
	}
}

func TestExtractFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		configData  string
		expected    []string
		expectError bool
	}{
		{
			name: "Valid config with repositories",
			configData: `
                module.exports = {
                    repositories: [
                        "group/project1",
                        "group/project2",
                        "group/subgroup/project3",
                    ],
                };
            `,
			expected:    []string{"group/project1", "group/project2", "group/subgroup/project3"},
			expectError: false,
		},
		{
			name: "Valid config with quotes",
			configData: `
			{
                "repositories": [
                        "group/project1",
                        "group/project2",
                        "group/subgroup/project3",
                    ]
                }
            `,
			expected:    []string{"group/project1", "group/project2", "group/subgroup/project3"},
			expectError: false,
		},
		{
			name: "Empty repositories section",
			configData: `
                module.exports = {
                    repositories: [
                    ],
                };
            `,
			expected:    nil,
			expectError: true,
		},
		{
			name: "No repositories section",
			configData: `
                module.exports = {
                    platform: "gitlab",
                };
            `,
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid config format",
			configData: `
                module.exports = {
                    repositories: [
                        "group/project1",
                        "group/project2",
                };
            `,
			expected:    nil,
			expectError: true,
		},
		{
			name: "Commented out project",
			configData: `
                module.exports = {
                    repositories: [
                        "group/project1",
                        "group/project2",
                        // "group/subgroup/project3",
                    ],
                };
            `,
			expected:    []string{"group/project1", "group/project2"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpfile, err := os.CreateTemp(t.TempDir(), "config.js")
			require.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.WriteString(tt.configData)
			require.NoError(t, err)
			err = tmpfile.Close()
			require.NoError(t, err)

			repositories, err := gl.ExtractFromFile(tmpfile.Name())

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, repositories)
			}
		})
	}
}
