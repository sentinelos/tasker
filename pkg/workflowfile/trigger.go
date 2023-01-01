package workflowfile

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// decodeTriggerBlock validates each part of the trigger block, building out a defined *Trigger
func decodeTriggerBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Trigger, hcl.Diagnostics) {
	trigger := &Trigger{
		Name:      block.Labels[0],
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description"},
			{Name: "conditions"},
		},
	})

	if !hclsyntax.ValidIdentifier(trigger.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid trigger name",
			Detail:   BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &trigger.Description))
	}

	if attr, exists := content.Attributes["conditions"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &trigger.Conditions))
	}

	return trigger, diags
}
