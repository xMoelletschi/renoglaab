package mergerequests

import (
	"github.com/sirupsen/logrus"
	"github.com/xMoelletschi/renoglaab/internal/config"
	gl "github.com/xMoelletschi/renoglaab/internal/gitlab"
)

func ReconcileProjectMergeRequests(config config.Config, repo string, client gl.ClientWrapper) {
	mrs := listProjectMergeRequests(config, repo, &client)
	for _, mr := range mrs {
		comment := config.Approve
		if config.AddComment {
			comment = config.Comment
		}

		err := createMergeRequestNote(repo, mr, comment, client.Client)

		if err != nil {
			logrus.WithError(err).Error("Failed to create Merge request note")
		}

		logrus.WithFields(logrus.Fields{"repository": repo, "mrID": mr}).Info("Approved MR")
	}
}
