package workflowfile

import (
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// decodeJobBlock validates each part of the job block, building out a defined *Job
func decodeJobBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Job, hcl.Diagnostics) {
	job := &Job{
		Name:      block.Labels[0],
		Trigger:   map[string]*Trigger{},
		Uses:      &Use{},
		Container: &Container{},
		Services:  map[string]*Container{},
		Steps:     map[string]*Step{},
		DeclRange: block.DefRange,
	}

	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "description"},
			{Name: "depends_on"},
			{Name: "timeout"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "trigger", LabelNames: []string{"name"}},
			{Type: "use", LabelNames: []string{"name"}},
			{Type: "container", LabelNames: []string{"image"}},
			{Type: "service", LabelNames: []string{"image"}},
			{Type: "step", LabelNames: []string{"name"}},
		},
	})

	if !hclsyntax.ValidIdentifier(job.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid job name",
			Detail:   BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &job.Description))
	}

	for _, blk := range content.Blocks.OfType("trigger") {
		trigger, triggerDiags := decodeTriggerBlock(blk, ctx)
		if triggerDiags.HasErrors() {
			diags = diags.Extend(triggerDiags)
		} else {
			if _, found := job.Trigger[trigger.Name]; found {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate trigger",
					Detail:   "Duplicate " + trigger.Name + " trigger definition found.",
					Subject:  &trigger.DeclRange,
					Context:  blk.DefRange.Ptr(),
				})
			} else {
				job.Trigger[trigger.Name] = trigger
			}
		}
	}

	if attr, exists := content.Attributes["depends_on"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &job.DependsOn))
	}

	for _, blk := range content.Blocks.OfType("use") {
		use, useDiags := decodeUseBlock(blk, ctx)
		if useDiags.HasErrors() {
			diags = diags.Extend(useDiags)
		} else {
			job.Uses = use
		}
	}

	for _, blk := range content.Blocks.OfType("container") {
		container, containerDiags := decodeContainerBlock(blk, ctx)
		if containerDiags.HasErrors() {
			diags = diags.Extend(containerDiags)
		} else {
			job.Container = container
		}
	}

	for _, blk := range content.Blocks.OfType("service") {
		service, serviceDiags := decodeContainerBlock(blk, ctx)
		if serviceDiags.HasErrors() {
			diags = diags.Extend(serviceDiags)
		} else {
			if _, found := job.Services[service.Image]; found {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate service",
					Detail:   "Duplicate " + service.Image + " service definition found.",
					Subject:  &service.DeclRange,
					Context:  blk.DefRange.Ptr(),
				})
			} else {
				job.Services[service.Image] = service
			}
		}
	}

	for _, blk := range content.Blocks.OfType("step") {
		step, stepDiags := decodeStepBlock(blk, ctx)
		if stepDiags.HasErrors() {
			diags = diags.Extend(stepDiags)
		} else {
			if _, found := job.Steps[step.Name]; found {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate step",
					Detail:   "Duplicate " + step.Name + " step definition found.",
					Subject:  &step.DeclRange,
					Context:  blk.DefRange.Ptr(),
				})
			} else {
				job.Steps[step.Name] = step
			}
		}
	}

	if attr, exists := content.Attributes["timeout"]; exists {
		var timeout string
		timeoutDiags := gohcl.DecodeExpression(attr.Expr, ctx, &timeout)
		diags = diags.Extend(timeoutDiags)

		if !timeoutDiags.HasErrors() {
			d, err := time.ParseDuration(timeout)
			if err != nil {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Invalid timeout",
					Detail:   fmt.Sprintf("Invalid timeout '%s', using default of %s", timeout, DefaultTimeout),
					Subject:  attr.Expr.StartRange().Ptr(),
					Context:  attr.Expr.Range().Ptr(),
				})

				job.Timeout = DefaultTimeout
			} else {
				job.Timeout = d
			}
		}
	}

	return job, diags
}
