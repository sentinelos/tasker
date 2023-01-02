package taskfile

import (
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// decodeTaskBlock validates each part of the task block, building out a defined *Task
func decodeTaskBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Task, hcl.Diagnostics) {
	task := &Task{
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

	if !hclsyntax.ValidIdentifier(task.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid task name",
			Detail:   BadIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &task.Description))
	}

	for _, blk := range content.Blocks.OfType("trigger") {
		trigger, triggerDiags := decodeTriggerBlock(blk, ctx)
		if triggerDiags.HasErrors() {
			return task, diags.Extend(triggerDiags)
		}

		if _, found := task.Trigger[trigger.Name]; found {
			return task, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate trigger",
				Detail:   "Duplicate " + trigger.Name + " trigger definition found.",
				Subject:  &trigger.DeclRange,
				Context:  blk.DefRange.Ptr(),
			})
		}

		task.Trigger[trigger.Name] = trigger
	}

	if attr, exists := content.Attributes["depends_on"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, ctx, &task.DependsOn))
	}

	for _, blk := range content.Blocks.OfType("use") {
		use, useDiags := decodeUseBlock(blk, ctx)
		if useDiags.HasErrors() {
			return task, diags.Extend(useDiags)
		}

		task.Uses = use
	}

	for _, blk := range content.Blocks.OfType("container") {
		container, containerDiags := decodeContainerBlock(blk, ctx)
		if containerDiags.HasErrors() {
			return task, diags.Extend(containerDiags)
		}

		task.Container = container
	}

	for _, blk := range content.Blocks.OfType("service") {
		service, serviceDiags := decodeContainerBlock(blk, ctx)
		if serviceDiags.HasErrors() {
			return task, diags.Extend(serviceDiags)
		}

		if _, found := task.Services[service.Image]; found {
			return task, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate service",
				Detail:   "Duplicate " + service.Image + " service definition found.",
				Subject:  &service.DeclRange,
				Context:  blk.DefRange.Ptr(),
			})
		}

		task.Services[service.Image] = service
	}

	for _, blk := range content.Blocks.OfType("step") {
		step, stepDiags := decodeStepBlock(blk, ctx)
		if stepDiags.HasErrors() {
			return task, diags.Extend(stepDiags)
		}

		if _, found := task.Steps[step.Name]; found {
			return task, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate step",
				Detail:   "Duplicate " + step.Name + " step definition found.",
				Subject:  &step.DeclRange,
				Context:  blk.DefRange.Ptr(),
			})
		}

		task.Steps[step.Name] = step
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

				task.Timeout = DefaultTimeout
			} else {
				task.Timeout = d
			}
		}
	}

	return task, diags
}
