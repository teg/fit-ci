package api

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"

	"github.com/palantir/go-githubapp/githubapp"
)

type checkSuiteHandler struct {
	githubapp.ClientCreator
}

func (h *checkSuiteHandler) Handles() []string {
	return []string{"check_suite"}
}

func (h *checkSuiteHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.CheckSuiteEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse check suite event payload")
	}

	repo := event.GetRepo()
	installationID := githubapp.GetInstallationIDFromEvent(&event)

	ctx, logger := githubapp.PrepareRepoContext(ctx, installationID, repo)

	logger.Debug().Msgf("Event action is %s", event.GetAction())
	if event.GetAction() != "requested" {
		return nil
	}

	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	headBranch := event.GetCheckSuite().GetHeadBranch()
	headSHA := event.GetCheckSuite().GetHeadSHA()

	logger.Debug().Msgf("Starting check run on %s/%s#%s", repoOwner, repoName, headBranch)

	options := github.CreateCheckRunOptions{
		Name:       "integration",
		HeadBranch: headBranch,
		HeadSHA:    headSHA,
	}

	_, _, err = client.Checks.CreateCheckRun(ctx, repoOwner, repoName, options)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create check run")
	}

	return nil
}
