package taskfile

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// decodeNotifyBlock validates each part of the notify block, building out a defined *Notify
func decodeNotifyBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Notify, hcl.Diagnostics) {
	notify := &Notify{
		Name:      block.Labels[0],
		Outputs:   map[string]*Output{},
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "output", LabelNames: []string{"name"}},
		},
	})

	if !hclsyntax.ValidIdentifier(notify.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid notify name",
			Detail:   BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &notify.Description))
	}

	for _, blk := range content.Blocks.OfType("output") {
		output, outputDiags := decodeOutputBlock(blk, ctx)
		if outputDiags.HasErrors() {
			return notify, diags.Extend(outputDiags)
		}

		if _, found := notify.Outputs[output.Name]; found {
			return notify, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate output",
				Detail:   "Duplicate " + output.Name + " output definition found.",
				Subject:  &output.DeclRange,
				Context:  blk.DefRange.Ptr(),
			})
		}

		notify.Outputs[output.Name] = output
	}

	return notify, diags
}
