package taskfile

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/sentinelos/tasker/internal/constants"
)

// decodeContainerBlock validates each part of the container block, building out a defined *Container
func decodeContainerBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Container, hcl.Diagnostics) {
	container := &Container{
		Image:      block.Labels[0],
		Credential: &Credential{},
		DeclRange:  block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description"},
			{Name: "environments"},
			{Name: "volumes"},
			{Name: "flags"},
			{Name: "command"},
			{Name: "args"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "credential"},
		},
	})

	if !hclsyntax.ValidIdentifier(container.Image) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid container image name",
			Detail:   constants.BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &container.Description))
	}

	if attr, exists := content.Attributes["environments"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &container.Environments))
	}

	if attr, exists := content.Attributes["volumes"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &container.Volumes))
	}

	if attr, exists := content.Attributes["flags"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &container.Flags))
	}

	if attr, exists := content.Attributes["command"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &container.Command))
	}

	if attr, exists := content.Attributes["args"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &container.Args))
	}

	for _, blk := range content.Blocks.OfType("credential") {
		credential, credentialDiags := decodeCredentialBlock(blk, ctx)
		if credentialDiags.HasErrors() {
			return container, diags.Extend(credentialDiags)
		}

		container.Credential = credential
	}

	return container, diags
}
