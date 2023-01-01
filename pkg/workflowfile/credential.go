package workflowfile

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// decodeCredentialBlock validates each part of the credential block, building out a defined *Credential
func decodeCredentialBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Credential, hcl.Diagnostics) {
	credential := &Credential{
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "username", Required: true},
			{Name: "password", Required: true},
		},
	})

	if attr, exists := content.Attributes["username"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &credential.Username))
	}

	if attr, exists := content.Attributes["password"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &credential.Password))
	}

	return credential, diags
}
