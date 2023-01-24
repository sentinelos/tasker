package taskfile

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/sentinelos/tasker/internal/constants"
)

// decodeNotifierBlock validates each part of the notifier block, building out a defined *Notifier
func decodeNotifierBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Notifier, hcl.Diagnostics) {
	notifier := &Notifier{
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

	if !hclsyntax.ValidIdentifier(notifier.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid notifier name",
			Detail:   constants.BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &notifier.Description))
	}

	for _, blk := range content.Blocks.OfType("output") {
		output, outputDiags := decodeOutputBlock(blk, ctx)
		if outputDiags.HasErrors() {
			return notifier, diags.Extend(outputDiags)
		}

		if _, found := notifier.Outputs[output.Name]; found {
			return notifier, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate output",
				Detail:   "Duplicate " + output.Name + " output definition found.",
				Subject:  &output.DeclRange,
				Context:  blk.DefRange.Ptr(),
			})
		}

		notifier.Outputs[output.Name] = output
	}

	return notifier, diags
}
