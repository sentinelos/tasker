package workflowfile

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// decodeInputBlock validates each part of the input block, building out a defined *Input
func decodeInputBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Input, hcl.Diagnostics) {
	input := &Input{
		Name:      block.Labels[0],
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description"},
			{Name: "value", Required: true},
			{Name: "sensitive"},
		},
	})

	if !hclsyntax.ValidIdentifier(input.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid input name",
			Detail:   BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &input.Description))
	}

	if attr, exists := content.Attributes["value"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &input.Value))
	}

	if attr, exists := content.Attributes["sensitive"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &input.Sensitive))
	}

	return input, diags
}
