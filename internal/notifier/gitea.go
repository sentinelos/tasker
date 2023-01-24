package notifier

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitea"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// NewGiteaNotifier returns a notifier that posts task steps statuses as status checks on a pull request.
func NewGiteaNotifier(namespace, name, ref string) (*Gitea, error) {
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
		GitNotifier: &GitNotifier{
			Namespace: namespace,
			Name:      name,
			Ref:       ref,
			client:    client,
			ctx:       context.Background(),
		},
	}, nil
}
