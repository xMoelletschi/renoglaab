//nolint:paralleltest,revive,unparam
package gitlab_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	gl "github.com/xMoelletschi/renoglaab/internal/gitlab"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func TestCreateGitLabClient(t *testing.T) {
	var fatalCalled bool

	logrus.StandardLogger().ExitFunc = func(int) { fatalCalled = true }

	tests := []struct {
		name          string
		gitlabToken   string
		gitlabBaseURL string
		mockError     error
		expectFatal   bool
	}{
		{
			name:          "Valid token and URL",
			gitlabToken:   "valid_token",
			gitlabBaseURL: "https://gitlab.com",
			mockError:     nil,
			expectFatal:   false,
		},
		{
			name:          "Empty token",
			gitlabToken:   "",
			gitlabBaseURL: "https://gitlab.com",
			mockError:     nil,
			expectFatal:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fatalCalled = false

			// Mock the gitlab.NewClient function
			newGitLabClient := func(_ string, options ...gitlab.ClientOptionFunc) (*gitlab.Client, error) {
				return nil, tt.mockError
			}

			// Temporarily replace the global newGitLabClient with the local mock
			originalNewGitLabClient := newGitLabClient
			defer func() { newGitLabClient = originalNewGitLabClient }()

			client, _ := gl.CreateGitLabClient(tt.gitlabToken, tt.gitlabBaseURL)

			if tt.expectFatal {
				assert.True(t, fatalCalled, "Expected fatal to be called")
			} else {
				assert.NotNil(t, client, "Expected client to be created")
				assert.False(t, fatalCalled, "Expected fatal not to be called")
			}
		})
	}
}
