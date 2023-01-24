package notifier

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type (
	// Notifier describes a hook for sending summarized results.
	Notifier interface {
		SetStatus(*scm.StatusInput) error
		Statuses(scm.ListOptions) ([]*scm.Status, error)
	}

	// GitNotifier describes a hook for sending summarized results to git gateway.
	GitNotifier struct {
		Namespace string
		Name      string
		Ref       string
		client    *scm.Client
		ctx       context.Context
	}

	// Bitbucket is a notifier that summarizes task steps statuses as Bitbucket statuses.
	Bitbucket struct {
		*GitNotifier
	}

	// Gitea is a notifier that summarizes task steps statuses as Gitea statuses.
	Gitea struct {
		*GitNotifier
	}

	// GitHub is a notifier that summarizes task steps statuses as GitHub statuses.
	GitHub struct {
		*GitNotifier
	}

	// GitLab is a notifier that summarizes task steps statuses as GitLab statuses.
	GitLab struct {
		*GitNotifier
	}

	// Noop is a notifier that does nothing.
	Noop struct{}
)
