package taskfile

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/sentinelos/tasker/internal/constants"
)

// decodeUseBlock validates each part of the use block, building out a defined *Use
func decodeUseBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Use, hcl.Diagnostics) {
	use := &Use{
		Name:      block.Labels[0],
		Inputs:    map[string]*Input{},
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "input", LabelNames: []string{"name"}},
		},
	})

	if !hclsyntax.ValidIdentifier(use.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid use name",
			Detail:   constants.BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &use.Description))
	}

	for _, blk := range content.Blocks.OfType("input") {
		input, inputDiags := decodeInputBlock(blk, ctx)
		if inputDiags.HasErrors() {
			return use, diags.Extend(inputDiags)
		}

		if _, found := use.Inputs[input.Name]; found {
			return use, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate input",
				Detail:   "Duplicate " + input.Name + " input definition found.",
				Subject:  &input.DeclRange,
				Context:  blk.DefRange.Ptr(),
			})
		}

		use.Inputs[input.Name] = input
	}

	return use, diags
}
