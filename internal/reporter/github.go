package reporter

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// GitHub is a reporter that summarizes workflow statuses as GitHub statuses.
type GitHub struct {
	*GitReporter
}

// NewGitHubReporter returns a reporter that posts workflow statuses as status checks on a pull request.
func NewGitHubReporter(namespace, name, ref string) (*GitHub, error) {
	uri, ok := os.LookupEnv("GITHUB_URI")
	if !ok {
		uri = "https://api.github.com"
	}

	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		return nil, errors.New("missing GITHUB_TOKEN")
	}

	client, err := github.New(uri)
	if err != nil {
		return nil, err
	}

	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&scm.Token{Token: token}),
		},
	}

	return &GitHub{
		GitReporter: &GitReporter{
			Namespace: namespace,
			Name:      name,
			Ref:       ref,
			client:    client,
			ctx:       context.Background(),
		},
	}, nil
}
