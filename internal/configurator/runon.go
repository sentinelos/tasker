package configurator

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/sentinelos/tasker/internal/constants"
	"github.com/sentinelos/tasker/internal/diagnostic/metadata"
	"golang.org/x/exp/slices"

	"github.com/containerd/containerd/platforms"
)

// decodeRunOnBlock validates each part of the runOn block, building out a defined *RunOn
func decodeRunOnBlock(block *hcl.Block) (*RunOn, hcl.Diagnostics) {
	workdir, _ := os.Getwd()
	hostMetadata := metadata.Get("host")

	runOn := &RunOn{
		Name:      block.Labels[0],
		Type:      block.Labels[1],
		Platform:  hostMetadata.GetLabel("platform"),
		User:      hostMetadata.GetLabel("user"),
		Shell:     hostMetadata.GetLabel("shell"),
		Workdir:   workdir,
		DeclRange: block.DefRange,
	}

	if !hclsyntax.ValidIdentifier(runOn.Name) {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Invalid run_on name",
			Detail:   constants.BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		}}
	}

	if !slices.Contains([]string{"host", "oci"}, runOn.Type) {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Invalid run_on type",
			Detail:   "Type \"" + runOn.Type + "\" seems incompatible, suitable values \"host\" or \"oci\"",
			Subject:  &block.LabelRanges[1],
		}}
	}

	attributes := []hcl.AttributeSchema{
		{Name: "description"},
		{Name: "environments"},
		{Name: "shell"},
	}

	if runOn.Type == "oci" {
		attributes = append(attributes, []hcl.AttributeSchema{
			{Name: "image", Required: true},
			{Name: "labels"},
			{Name: "cpu"},
			{Name: "memory"},
			{Name: "platform"},
			{Name: "user"},
			{Name: "workdir"},
		}...)
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: attributes,
	})

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Description))
	}

	if attr, exists := content.Attributes["environments"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Environments))
	}

	if runOn.Type == "oci" {
		if attr, exists := content.Attributes["image"]; exists {
			var image string
			diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &image))

		}

		if attr, exists := content.Attributes["labels"]; exists {
			diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Labels))
		}

		if attr, exists := content.Attributes["cpu"]; exists {
			diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Cpu))
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
					Detail:   "Platform \"" + runOn.Platform + "\" seems incompatible with the host platform \"" + hostMetadata.GetLabel("platform") + "\"",
					Subject:  attr.Expr.StartRange().Ptr(),
					Context:  attr.Expr.Range().Ptr(),
				})
			}
		}

		if attr, exists := content.Attributes["user"]; exists {
			diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.User))
		}

		if attr, exists := content.Attributes["workdir"]; exists {
			diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &runOn.Workdir))
		}
	}

	return runOn, diags
}
