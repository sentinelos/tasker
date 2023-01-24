package notifier

import (
	"fmt"

	"github.com/drone/go-scm/scm"
)

// SetStatus creates a new commit status.
func (gr *GitNotifier) SetStatus(input *scm.StatusInput) error {
	_, _, err := gr.client.Repositories.CreateStatus(gr.ctx, gr.repository(), gr.Ref, input)
	return err
}

// Statuses returns a list of commit statuses.
func (gr *GitNotifier) Statuses(options scm.ListOptions) ([]*scm.Status, error) {
	statuses, _, err := gr.client.Repositories.ListStatus(gr.ctx, gr.repository(), gr.Ref, options)
	return statuses, err
}

func (gr *GitNotifier) repository() string {
	return fmt.Sprintf("%s/%s", gr.Namespace, gr.Name)
}
