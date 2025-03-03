package gitlab

import (
	"github.com/sirupsen/logrus"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Ensure ClientWrapper implements Client interface.
var _ Client = (*ClientWrapper)(nil)

// ClientWrapper wraps the gitlab.Client to implement the Client interface.
type ClientWrapper struct {
	Client *gitlab.Client
}

// Client defines the methods used by ListProjectMergeRequests and pipelineSucceeded.
type Client interface {
	ListProjectMergeRequests(
		repo string, opts *gitlab.ListProjectMergeRequestsOptions,
	) ([]*gitlab.BasicMergeRequest, *gitlab.Response, error)
	ListProjectPipelines(
		repo string, opts *gitlab.ListProjectPipelinesOptions,
	) ([]*gitlab.PipelineInfo, *gitlab.Response, error)
	GetPipeline(
		repo string, pipelineID int,
	) (*gitlab.Pipeline, *gitlab.Response, error)
}

// ListProjectMergeRequests fetches the merge requests for a given repository.
func (w *ClientWrapper) ListProjectMergeRequests(
	repo string, opts *gitlab.ListProjectMergeRequestsOptions,
) ([]*gitlab.BasicMergeRequest, *gitlab.Response, error) {
	logrus.WithFields(logrus.Fields{
		"repo": repo,
	}).Debug("Fetching merge requests")

	return w.Client.MergeRequests.ListProjectMergeRequests(repo, opts)
}

// ListProjectPipelines fetches the pipelines for a given repository.
func (w *ClientWrapper) ListProjectPipelines(
	repo string, opts *gitlab.ListProjectPipelinesOptions,
) ([]*gitlab.PipelineInfo, *gitlab.Response, error) {
	logrus.WithFields(logrus.Fields{
		"repo": repo,
	}).Debug("Fetching pipelines")

	return w.Client.Pipelines.ListProjectPipelines(repo, opts)
}

// GetPipeline fetches a specific pipeline by its ID for a given repository.
func (w *ClientWrapper) GetPipeline(
	repo string, pipelineID int,
) (*gitlab.Pipeline, *gitlab.Response, error) {
	logrus.WithFields(logrus.Fields{
		"repo":       repo,
		"pipelineID": pipelineID,
	}).Debug("Fetching pipeline")

	return w.Client.Pipelines.GetPipeline(repo, pipelineID)
}

// CreateGitLabClient initializes a new GitLab client.
func CreateGitLabClient(gitlabToken string, gitlabBaseURL string) (*ClientWrapper, error) {
	if gitlabToken == "" {
		logrus.Fatal("GITLAB_API_TOKEN must be set")
	}

	client, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabBaseURL))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create GitLab client")
	}

	logrus.Debug("GitLab client successfully initialized")

	return &ClientWrapper{Client: client}, nil
}
