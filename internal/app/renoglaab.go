package app

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/xMoelletschi/renoglaab/internal/config"
	gl "github.com/xMoelletschi/renoglaab/internal/gitlab"
	"github.com/xMoelletschi/renoglaab/internal/mergerequests"
)

const workerCount = 5

// Run is the main entry point for executing the application logic.
// It performs the following steps:
// 1. Loads the configuration from the config file.
// 2. Extracts the list of repositories from the configuration.
// 3. Creates a GitLab client using the provided API token and URL.
// 4. Iterates over each repository and reconciles the merge requests.
func Run() error {
	cfg := config.NewConfig()

	repositories, err := gl.ExtractRepositories(cfg.ConfigPath)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to extract repositories from config path: %s", cfg.ConfigPath)

		return err
	}

	gitLabClient, err := gl.CreateGitLabClient(cfg.GitLabAPIToken, cfg.GitLabURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to create GitLab client")

		return err
	}

	repoChan := make(chan string, len(repositories))

	var wg sync.WaitGroup

	for i := range make([]struct{}, workerCount) {
		wg.Add(1)

		go func(_ int) {
			defer wg.Done()

			for repo := range repoChan {
				mergerequests.ReconcileProjectMergeRequests(*cfg, repo, *gitLabClient)
			}
		}(i)
	}

	for _, repo := range repositories {
		repoChan <- repo
	}

	close(repoChan)

	wg.Wait()

	return nil
}
