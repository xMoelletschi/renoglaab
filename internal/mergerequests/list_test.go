//nolint:lll,funlen,err113
package mergerequests

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xMoelletschi/renoglaab/internal/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// MockGitLabClient is a mock implementation of GitLabClient.
type MockGitLabClient struct {
	mock.Mock
}

func (m *MockGitLabClient) ListProjectMergeRequests(repo string, opts *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	args := m.Called(repo, opts)
	mrs, ok := args.Get(0).([]*gitlab.MergeRequest)

	if !ok {
		return nil, nil, errors.New("type assertion to []*gitlab.MergeRequest failed")
	}

	return mrs, nil, args.Error(1)
}

func (m *MockGitLabClient) ListProjectPipelines(repo string, opts *gitlab.ListProjectPipelinesOptions) ([]*gitlab.PipelineInfo, *gitlab.Response, error) {
	args := m.Called(repo, opts)
	pipelines, ok := args.Get(0).([]*gitlab.PipelineInfo)

	if !ok {
		return nil, nil, errors.New("type assertion to []*gitlab.PipelineInfo failed")
	}

	return pipelines, nil, args.Error(1)
}

func (m *MockGitLabClient) GetPipeline(repo string, pipelineID int) (*gitlab.Pipeline, *gitlab.Response, error) {
	args := m.Called(repo, pipelineID)
	pipeline, ok := args.Get(0).(*gitlab.Pipeline)

	if !ok {
		return nil, nil, errors.New("type assertion to *gitlab.Pipeline failed")
	}

	return pipeline, nil, args.Error(1)
}

func TestListProjectMergeRequests(t *testing.T) {
	t.Parallel()

	repo := "test/repo"
	config := config.Config{
		Labels:                    []string{"ops", "renovate"},
		AllowedBranchRegex:        regexp.MustCompile(`^feature/branch\d+$`),
		FilterByBranch:            true,
		FilterBySucceededPipeline: true,
	}

	tests := []struct {
		name      string
		mrs       []*gitlab.MergeRequest
		pipelines []*gitlab.PipelineInfo
		pipeline  *gitlab.Pipeline
		expectIDs []int
		listErr   error
		pipeErr   error
		getErr    error
	}{
		{
			name: "Valid MR and successful pipeline",
			mrs: []*gitlab.MergeRequest{
				{IID: 1, SourceBranch: "feature/branch1", Title: "Test MR"},
			},
			pipelines: []*gitlab.PipelineInfo{{ID: 100}},
			pipeline:  &gitlab.Pipeline{Status: "success", DetailedStatus: &gitlab.DetailedStatus{Icon: "status_success"}},
			expectIDs: []int{1},
		},
		{
			name: "MR with failing pipeline",
			mrs: []*gitlab.MergeRequest{
				{IID: 2, SourceBranch: "feature/branch2", Title: "Test MR 2", Pipeline: &gitlab.PipelineInfo{Status: "failed"}},
			},
			pipelines: []*gitlab.PipelineInfo{{ID: 102}},
			pipeline:  &gitlab.Pipeline{Status: "failed", DetailedStatus: &gitlab.DetailedStatus{Icon: "failed"}},
			expectIDs: nil,
		},
		{
			name: "MR with no pipelines",
			mrs: []*gitlab.MergeRequest{
				{IID: 3, SourceBranch: "feature/branch3", Title: "Test MR 3"},
			},
			pipelines: []*gitlab.PipelineInfo{},
			expectIDs: nil,
		},
		{
			name: "MR with a branch that does not match regex",
			mrs: []*gitlab.MergeRequest{
				{IID: 4, SourceBranch: "hotfix/branch4", Title: "Hotfix MR"},
			},
			expectIDs: nil,
		},
		{
			name:      "GitLab API error on listing MRs",
			mrs:       nil,
			listErr:   errors.New("GitLab API error"),
			expectIDs: nil,
		},
		{
			name: "Error listing pipelines",
			mrs: []*gitlab.MergeRequest{
				{IID: 5, SourceBranch: "feature/branch5", Title: "Test MR 5"},
			},
			pipelines: nil,
			expectIDs: nil,
			pipeErr:   errors.New("Failed to list pipelines"),
		},
		{
			name: "Error getting pipeline details",
			mrs: []*gitlab.MergeRequest{
				{IID: 6, SourceBranch: "feature/branch6", Title: "Test MR 6"},
			},
			pipelines: []*gitlab.PipelineInfo{{ID: 103}},
			expectIDs: nil,
			getErr:    errors.New("Failed to get pipeline details"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(MockGitLabClient)
			mockClient.ExpectedCalls = nil
			mockClient.Calls = nil

			mockClient.On("ListProjectMergeRequests", repo, mock.Anything).Return(tt.mrs, tt.listErr).Once()

			if tt.mrs != nil {
				mockClient.On("ListProjectPipelines", repo, mock.Anything).Return(tt.pipelines, tt.pipeErr).Maybe()

				if len(tt.pipelines) > 0 {
					mockClient.On("GetPipeline", repo, tt.pipelines[0].ID).Return(tt.pipeline, tt.getErr).Maybe()
				}
			}

			result := listProjectMergeRequests(config, repo, mockClient)
			assert.Equal(t, tt.expectIDs, result)
		})
	}
}
