package taskfile

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// decodeOutputBlock validates each part of the output block, building out a defined *Output
func decodeOutputBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Output, hcl.Diagnostics) {
	output := &Output{
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

	if !hclsyntax.ValidIdentifier(output.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid output name",
			Detail:   BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &output.Description))
	}

	if attr, exists := content.Attributes["value"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &output.Value))
	}

	if attr, exists := content.Attributes["sensitive"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &output.Sensitive))
	}

	return output, diags
}
