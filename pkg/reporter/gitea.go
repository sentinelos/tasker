package reporter

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitea"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// Gitea is a reporter that summarizes workflow statuses as Gitea statuses.
type Gitea struct {
	*GitReporter
}

// NewGiteaReporter returns a reporter that posts workflow statuses as status checks on a pull request.
func NewGiteaReporter(namespace, name, ref string) (*Gitea, error) {
	uri, ok := os.LookupEnv("GITEA_URI")
	if !ok {
		return nil, errors.New("missing GITEA_URI")
	}

	token, ok := os.LookupEnv("GITEA_TOKEN")
	if !ok {
		return nil, errors.New("missing GITEA_TOKEN")
	}

	client, err := gitea.New(uri)
	if err != nil {
		return nil, err
	}

	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&scm.Token{Token: token}),
		},
	}

	return &Gitea{
		GitReporter: &GitReporter{
			Namespace: namespace,
			Name:      name,
			Ref:       ref,
			client:    client,
			ctx:       context.Background(),
		},
	}, nil
}
