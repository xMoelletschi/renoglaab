package mergerequests

import (
	"github.com/sirupsen/logrus"
	"github.com/xMoelletschi/renoglaab/internal/config"
	gl "github.com/xMoelletschi/renoglaab/internal/gitlab"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

const stateOpen string = "opened"

// listProjectMergeRequests fetches MRs and filters by branch regex and pipeline status.
func listProjectMergeRequests(config config.Config, repo string, client gl.Client) []int {
	logrus.WithField("repository", repo).Debug("Listing merge requests")

	options := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.Ptr(stateOpen),
	}

	if config.FilterByAuthorUsername {
		options.AuthorUsername = &config.AuthorUsername
	}

	if config.FilterByLabels {
		labels := gitlab.LabelOptions(config.Labels)
		options.Labels = &labels
	}

	mrs, _, err := client.ListProjectMergeRequests(repo, options)
	if err != nil {
		logrus.WithError(err).WithField("repository", repo).Error("Failed to list MRs")
	}

	var mrIIDs []int

	for _, mr := range mrs {
		if shouldProcessMR(repo, mr, config, client) {
			mrIIDs = append(mrIIDs, mr.IID)
		}
	}

	return mrIIDs
}

func shouldProcessMR(repo string, mr *gitlab.MergeRequest, config config.Config, client gl.Client) bool {
	logrus.WithFields(logrus.Fields{
		"repository": repo, "mr_id": mr.IID, "branch": mr.SourceBranch, "title": mr.Title,
	}).Debug("Checking")

	if config.FilterByBranch {
		if !config.AllowedBranchRegex.MatchString(mr.SourceBranch) {
			logrus.WithFields(logrus.Fields{
				"repository": repo, "mr_id": mr.IID, "branch": mr.SourceBranch, "title": mr.Title,
			}).Debug("Branch does not match allowed regex")

			return false
		}

		logrus.WithFields(logrus.Fields{
			"repository": repo, "mr_id": mr.IID, "branch": mr.SourceBranch, "title": mr.Title,
		}).Debug("Branch matches allowed regex")
	}

	if config.FilterBySucceededPipeline {
		if !pipelineSucceeded(config, repo, mr.SourceBranch, client) {
			logrus.WithFields(logrus.Fields{
				"repository": repo, "mr_id": mr.IID, "branch": mr.SourceBranch, "title": mr.Title,
			}).Debug("Pipeline failed for MR")

			return false
		}

		logrus.WithFields(logrus.Fields{
			"repository": repo, "mr_id": mr.IID, "branch": mr.SourceBranch, "title": mr.Title,
		}).Debug("Pipeline succeeded for MR")
	}

	return true
}
