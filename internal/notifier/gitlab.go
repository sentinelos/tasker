package notifier

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitlab"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// NewGitLabNotifier returns a notifier that posts task steps statuses as status checks on a pull request.
func NewGitLabNotifier(namespace, name, ref string) (*GitLab, error) {
	uri, ok := os.LookupEnv("GITLAB_URI")
	if !ok {
		uri = "https://gitlab.com"
	}

	token, ok := os.LookupEnv("GITLAB_TOKEN")
	if !ok {
		return nil, errors.New("missing GITLAB_TOKEN")
	}

	client, err := gitlab.New(uri)
	if err != nil {
		return nil, err
	}

	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&scm.Token{Token: token}),
		},
	}

	return &GitLab{
		GitNotifier: &GitNotifier{
			Namespace: namespace,
			Name:      name,
			Ref:       ref,
			client:    client,
			ctx:       context.Background(),
		},
	}, nil
}
