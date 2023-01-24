package configurator

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/sentinelos/tasker/internal/constants"
	"github.com/sentinelos/tasker/internal/diagnostic/metadata"

	"github.com/containerd/containerd/platforms"
)

// decodeRunOnBlock validates each part of the runOn block, building out a defined *RunOn
func decodeRunOnBlock(block *hcl.Block) (*RunOn, hcl.Diagnostics) {
	runOn := &RunOn{
		Name:      block.Labels[0],
		Image:     block.Labels[1],
		Platform:  platforms.DefaultString(),
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description"},
			{Name: "environments"},
			{Name: "cpu"},
			{Name: "memory"},
			{Name: "platform"},
			{Name: "user"},
			{Name: "shell", Required: true},
			{Name: "workdir"},
		},
	})

	if !hclsyntax.ValidIdentifier(runOn.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid run_on name",
			Detail:   constants.BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if !hclsyntax.ValidIdentifier(runOn.Image) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid run_on image",
			Detail:   constants.BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Description))
	}

	if attr, exists := content.Attributes["environments"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Environments))
	}

	if attr, exists := content.Attributes["cpu"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Cpu))
	}

	if attr, exists := content.Attributes["memory"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Memory))
	}

	if attr, exists := content.Attributes["platform"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Platform))

		if _, err := platforms.Parse(runOn.Platform); err != nil {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid run_on platform",
				Detail:   "Platform \"" + runOn.Platform + "\" seems incompatible with the host platform \"" + metadata.Get("host").GetLabel("platform") + "\"",
				Subject:  &block.LabelRanges[0],
			})
		}
	}

	if attr, exists := content.Attributes["user"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.User))
	}

	if attr, exists := content.Attributes["shell"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Shell))
	}

	if attr, exists := content.Attributes["workdir"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Workdir))
	}

	return runOn, diags
}
