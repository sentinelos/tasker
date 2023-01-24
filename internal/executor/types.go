package executor

import (
	"context"

	"github.com/hashicorp/hcl/v2"
)

type (
	Status uint

	ContainerStatus uint

	Context struct {
		Ctx *hcl.EvalContext
	}

	Runner interface {
		pre() Executor
		main() Executor
		post() Executor
	}

	Task interface {
		pre() Executor
		main() Executor
		post() Executor
	}

	Step interface {
		pre() Executor
		main() Executor
		post() Executor
	}

	// Executor define contract for the steps of a task
	Executor func(ctx context.Context) error
)
