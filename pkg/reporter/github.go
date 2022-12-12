package reporter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v41/github"
)

// GitHub is a reporter that summarizes policy statuses as GitHub statuses.
type GitHub struct {
	token string
	owner string
	repo  string
	sha   string
}

// NewGitHubReporter returns a reporter that posts policy checks as
// status checks on a pull request.
func NewGitHubReporter() (*GitHub, error) {
	token, ok := os.LookupEnv("INPUT_TOKEN")
	if !ok {
		return nil, errors.New("missing INPUT_TOKEN")
	}

	eventPath, ok := os.LookupEnv("GITHUB_EVENT_PATH")
	if !ok {
		return nil, errors.New("GITHUB_EVENT_PATH is not set")
	}

	data, err := os.ReadFile(eventPath)
	if err != nil {
		return nil, err
	}

	pullRequestEvent := &github.PullRequestEvent{}

	if err = json.Unmarshal(data, pullRequestEvent); err != nil {
		return nil, err
	}

	gh := &GitHub{
		token: token,
		owner: pullRequestEvent.GetRepo().GetOwner().GetLogin(),
		repo:  pullRequestEvent.GetRepo().GetName(),
		sha:   pullRequestEvent.GetPullRequest().GetHead().GetSHA(),
	}

	return gh, nil
}

// SetStatus sets the status of a GitHub check.
//
// Valid statuses are "error", "failure", "pending", "success".
func (gh *GitHub) SetStatus(state, policy, check, message string) error {
	if gh.token == "" {
		return errors.New("no token")
	}

	statusCheckContext := strings.ReplaceAll(strings.ToLower(path.Join("ensurer", policy, check)), " ", "-")
	repoStatus := &github.RepoStatus{}
	repoStatus.Context = &statusCheckContext
	repoStatus.Description = &message
	repoStatus.State = &state

	http.DefaultClient.Transport = roundTripper{gh.token}
	githubClient := github.NewClient(http.DefaultClient)

	_, _, err := githubClient.Repositories.CreateStatus(context.Background(), gh.owner, gh.repo, gh.sha, repoStatus)
	if err != nil {
		return err
	}

	return nil
}

type roundTripper struct {
	accessToken string
}

// RoundTrip implements the net/http.RoundTripper interface.
func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt.accessToken))

	return http.DefaultTransport.RoundTrip(r)
}
