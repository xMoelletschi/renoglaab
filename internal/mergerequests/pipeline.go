package mergerequests

import (
	"github.com/sirupsen/logrus"
	"github.com/xMoelletschi/renoglaab/internal/config"
	gl "github.com/xMoelletschi/renoglaab/internal/gitlab"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// pipelineSucceeded checks if the latest pipeline for a branch succeeded without warnings.
func pipelineSucceeded(config config.Config, repo, branch string, client gl.Client) bool {
	logrus.WithFields(logrus.Fields{
		"repository": repo,
		"branch":     branch,
	}).Debug("Checking pipeline status for branch")

	pipelines, err := listPipelines(client, repo, branch)
	if err != nil || len(pipelines) == 0 {
		return false
	}

	latestPipeline := pipelines[0]
	pipeline, err := getPipeline(client, repo, latestPipeline.ID)

	if err != nil {
		return false
	}

	return checkPipelineStatus(config, repo, latestPipeline.ID, pipeline)
}

func listPipelines(client gl.Client, repo, branch string) ([]*gitlab.PipelineInfo, error) {
	pipelines, _, err := client.ListProjectPipelines(repo, &gitlab.ListProjectPipelinesOptions{
		Ref: &branch, // Filter by branch
	})
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"repository": repo,
			"branch":     branch,
		}).Error("Failed to list pipelines for branch")

		return nil, err
	}

	if len(pipelines) == 0 {
		logrus.WithFields(logrus.Fields{
			"repository": repo,
			"branch":     branch,
		}).Warn("No pipelines found for branch")
	}

	return pipelines, nil
}

func getPipeline(client gl.Client, repo string, pipelineID int) (*gitlab.Pipeline, error) {
	pipeline, _, err := client.GetPipeline(repo, pipelineID)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"repository":  repo,
			"pipeline_id": pipelineID,
		}).Error("Failed to get detailed pipeline info")

		return nil, err
	}

	return pipeline, nil
}

func checkPipelineStatus(config config.Config, repo string, pipelineID int, pipeline *gitlab.Pipeline) bool {
	if pipeline.Status != "success" {
		logrus.WithFields(logrus.Fields{
			"repository":      repo,
			"pipeline_id":     pipelineID,
			"pipeline_status": pipeline.Status,
			"pipeline_icon":   pipeline.DetailedStatus.Icon,
		}).Warn("Pipeline did not succeed")

		return false
	}

	if config.FilterByPipelineWithoutWarnings && pipeline.DetailedStatus.Icon != "status_success" {
		logrus.WithFields(logrus.Fields{
			"repository":    repo,
			"pipeline_id":   pipelineID,
			"pipeline_icon": pipeline.DetailedStatus.Icon,
		}).Warn("Pipeline succeeded but has warnings")

		return false
	}

	logrus.WithFields(logrus.Fields{
		"repository":  repo,
		"pipeline_id": pipelineID,
	}).Debug("Pipeline succeeded")

	return true
}
