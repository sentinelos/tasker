package taskfile

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/sentinelos/tasker/internal/constants"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/gocty"
)

// decodeVariableBlock validates each part of the variable block, building out a defined *Variable
func decodeVariableBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Variable, hcl.Diagnostics) {
	variable := &Variable{
		Name:      block.Labels[0],
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description"},
			{Name: "type"},
			{Name: "default"},
			{Name: "sensitive"},
		},
	})

	if !hclsyntax.ValidIdentifier(variable.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid variable name",
			Detail:   constants.BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &variable.Description))
	}

	varType := cty.String
	if attr, exists := content.Attributes["type"]; exists {
		vType, typeDiags := typeexpr.Type(attr.Expr)
		if typeDiags.HasErrors() {
			return variable, diags.Extend(typeDiags)
		}

		varType = vType
	}

	if attr, exists := content.Attributes["default"]; exists {
		srcVal, srcDiags := attr.Expr.Value(ctx)
		if srcDiags.HasErrors() {
			return variable, diags.Extend(srcDiags)
		}

		value, err := convert.Convert(srcVal, varType)
		if err != nil {
			return variable, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsuitable value type",
				Detail:   fmt.Sprintf("Unsuitable value: %s", err.Error()),
				Subject:  attr.Expr.StartRange().Ptr(),
				Context:  attr.Expr.Range().Ptr(),
			})
		}

		variable.Value = value
	} else {
		if envVar, err := gocty.ToCtyValue(os.Getenv(constants.VarEnvPrefix+variable.Name), varType); err == nil {
			variable.Value = envVar
		}
	}

	if attr, exists := content.Attributes["sensitive"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &variable.Sensitive))
	}

	return variable, diags
}
