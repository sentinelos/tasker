package taskfile

import (
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/sentinelos/tasker/internal/constants"
)

// decodeStepBlock validates each part of the step block, building out a defined *Step
func decodeStepBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Step, hcl.Diagnostics) {
	step := &Step{
		Name:      block.Labels[0],
		If:        true,
		Uses:      &Use{},
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description"},
			{Name: "if"},
			{Name: "on_failure"},
			{Name: "shell"},
			{Name: "workdir"},
			{Name: "run", Required: true},
			{Name: "timeout"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "use", LabelNames: []string{"name"}},
		},
	})

	if !hclsyntax.ValidIdentifier(step.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid step name",
			Detail:   constants.BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &step.Description))
	}

	if attr, exists := content.Attributes["if"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &step.If))
	}

	if attr, exists := content.Attributes["on_failure"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &step.OnFailure))
	}

	for _, blk := range content.Blocks.OfType("use") {
		use, useDiags := decodeUseBlock(blk, ctx)
		if useDiags.HasErrors() {
			return step, diags.Extend(useDiags)
		}

		step.Uses = use
	}

	if attr, exists := content.Attributes["shell"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &step.Shell))
	}

	if attr, exists := content.Attributes["workdir"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &step.Workdir))
	}

	if attr, exists := content.Attributes["run"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &step.Run))
	}

	if attr, exists := content.Attributes["timeout"]; exists {
		var timeout string
		timeoutDiags := gohcl.DecodeExpression(attr.Expr, ctx, &timeout)
		diags = diags.Extend(timeoutDiags)

		if !timeoutDiags.HasErrors() {
			d, err := time.ParseDuration(timeout)
			if err != nil {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Invalid timeout",
					Detail:   fmt.Sprintf("Invalid timeout '%s', using default of %s", timeout, constants.DefaultTimeout),
					Subject:  attr.Expr.StartRange().Ptr(),
					Context:  attr.Expr.Range().Ptr(),
				})

				step.Timeout = constants.DefaultTimeout
			} else {
				step.Timeout = d
			}
		}
	}

	return step, diags
}
