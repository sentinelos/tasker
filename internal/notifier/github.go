package notifier

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// NewGitHubNotifier returns a notifier that posts task steps statuses as status checks on a pull request.
func NewGitHubNotifier(namespace, name, ref string) (*GitHub, error) {
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
		GitNotifier: &GitNotifier{
			Namespace: namespace,
			Name:      name,
			Ref:       ref,
			client:    client,
			ctx:       context.Background(),
		},
	}, nil
}
