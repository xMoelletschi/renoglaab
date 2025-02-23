//nolint:err113,funlen,paralleltest
package mergerequests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xMoelletschi/renoglaab/internal/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func TestPipelineSucceeded(t *testing.T) {
	repo := "test/repo"
	branch := "feature/branch1"
	config := config.Config{
		FilterByPipelineWithoutWarnings: true,
	}

	tests := []struct {
		name      string
		pipelines []*gitlab.PipelineInfo
		pipeline  *gitlab.Pipeline
		listErr   error
		getErr    error
		expected  bool
	}{
		{
			name:      "No pipelines found",
			pipelines: []*gitlab.PipelineInfo{},
			expected:  false,
		},
		{
			name:     "Failed to list pipelines",
			listErr:  errors.New("Failed to list pipelines"),
			expected: false,
		},
		{
			name:      "Failed to get detailed pipeline info",
			pipelines: []*gitlab.PipelineInfo{{ID: 100}},
			getErr:    errors.New("Failed to get pipeline details"),
			expected:  false,
		},
		{
			name:      "Pipeline did not succeed",
			pipelines: []*gitlab.PipelineInfo{{ID: 100}},
			pipeline:  &gitlab.Pipeline{Status: "failed", DetailedStatus: &gitlab.DetailedStatus{Icon: "failed"}},
			expected:  false,
		},
		{
			name:      "Pipeline succeeded but has warnings",
			pipelines: []*gitlab.PipelineInfo{{ID: 100}},
			pipeline:  &gitlab.Pipeline{Status: "success", DetailedStatus: &gitlab.DetailedStatus{Icon: "status_warning"}},
			expected:  false,
		},
		{
			name:      "Pipeline succeeded without warnings",
			pipelines: []*gitlab.PipelineInfo{{ID: 100}},
			pipeline:  &gitlab.Pipeline{Status: "success", DetailedStatus: &gitlab.DetailedStatus{Icon: "status_success"}},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockGitLabClient)
			mockClient.ExpectedCalls = nil
			mockClient.Calls = nil

			mockClient.On("ListProjectPipelines", repo, mock.Anything).Return(tt.pipelines, tt.listErr).Once()

			if len(tt.pipelines) > 0 {
				mockClient.On("GetPipeline", repo, tt.pipelines[0].ID).Return(tt.pipeline, tt.getErr).Once()
			}

			result := pipelineSucceeded(config, repo, branch, mockClient)
			assert.Equal(t, tt.expected, result)
		})
	}
}
