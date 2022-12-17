// Package reporter provides check result reporting.
package reporter

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

// Reporter describes a hook for sending summarized results.
type Reporter interface {
	SetStatus(*scm.StatusInput) error
	Statuses(scm.ListOptions) ([]*scm.Status, error)
}

// GitReporter describes a hook for sending summarized results to git.
type GitReporter struct {
	Namespace string
	Name      string
	Ref       string
	client    *scm.Client
	ctx       context.Context
}

// SetStatus creates a new commit status.
func (gr *GitReporter) SetStatus(input *scm.StatusInput) error {
	_, _, err := gr.client.Repositories.CreateStatus(gr.ctx, gr.repository(), gr.Ref, input)
	return err
}

// Statuses returns a list of commit statuses.
func (gr *GitReporter) Statuses(options scm.ListOptions) ([]*scm.Status, error) {
	statuses, _, err := gr.client.Repositories.ListStatus(gr.ctx, gr.repository(), gr.Ref, options)
	return statuses, err
}

func (gr *GitReporter) repository() string {
	return fmt.Sprintf("%s/%s", gr.Namespace, gr.Name)
}
