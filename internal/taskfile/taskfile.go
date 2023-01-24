// Package taskfile defines file specification.
package taskfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sentinelos/tasker/internal/executor"
	"github.com/zclconf/go-cty/cty"
)

// LoadTaskFile is a wrapper around DecodeTaskFile that first reads the given filename from disk and parses.
func LoadTaskFile(filename string, context *executor.Context) (*TaskFile, hcl.Diagnostics) {
	src, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Task file not found",
				Detail:   fmt.Sprintf("The Task file %s does not exist.", filename),
				Subject:  &hcl.Range{Filename: filename},
			}}
		}
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read Task file",
			Detail:   fmt.Sprintf("The Task file %s could not be read, %s.", filename, err),
			Subject:  &hcl.Range{Filename: filename},
		}}
	}

	if len(src) == 0 {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Empty Task file",
			Detail:   fmt.Sprintf("The Task file %s is empty.", filename),
			Subject:  &hcl.Range{Filename: filename},
		}}
	}

	var (
		file  *hcl.File
		diags hcl.Diagnostics
	)

	parser := hclparse.NewParser()
	switch suffix := strings.ToLower(filepath.Ext(filename)); suffix {
	case ".tf":
		file, diags = parser.ParseHCL(src, filename)
	case ".tf.json":
		file, diags = parser.ParseJSON(src, filename)
	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported Task file format",
			Detail:   fmt.Sprintf("The Task file %s could not be read, unrecognized format suffix %s.", filename, suffix),
			Subject:  &hcl.Range{Filename: filename},
		})
	}

	if diags.HasErrors() {
		return &TaskFile{
			Filename: filename,
			Source:   file,
		}, diags
	}

	return DecodeTaskFile(filename, file, context)
}

// DecodeTaskFile decodes and evaluates expressions in the given TaskFile source code.
//
// The "filename" argument provides TaskFile filename
// The "file" argument provides parsed TaskFile
// The "context" argument provides variables and functions for use during expression evaluation.
func DecodeTaskFile(filename string, file *hcl.File, context *executor.Context) (*TaskFile, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	taskFile := &TaskFile{
		Filename:  filename,
		Source:    file,
		Notifiers: map[string]*Notifier{},
		Tasks:     map[string]*Task{},
	}

	content, contentDiags := file.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "name", Required: true},
			{Name: "description", Required: true},
			{Name: "runs_on", Required: true},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "notifier", LabelNames: []string{"name"}},
			{Type: "task", LabelNames: []string{"name"}},
		},
	})

	diags = diags.Extend(contentDiags)

	if attr, exists := content.Attributes["name"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &taskFile.Name))
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &taskFile.Description))
	}

	if attr, exists := content.Attributes["runs_on"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &taskFile.RunsOn))
	}

	variables := map[string]cty.Value{}
	for _, block := range content.Blocks.OfType("variable") {
		variable, varDiags := decodeVariableBlock(block, nil)
		if varDiags.HasErrors() {
			return taskFile, diags.Extend(varDiags)
		}

		if _, found := variables[variable.Name]; found {
			return taskFile, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate variable",
				Detail:   "Duplicate " + variable.Name + " variable definition found.",
				Subject:  &variable.DeclRange,
				Context:  block.DefRange.Ptr(),
			})
		}

		variables[variable.Name] = variable.Value
	}

	context.AddVariables(variables)

	ctx := context.Ctx
	for _, block := range content.Blocks.OfType("notifier") {
		notifier, notifierDiags := decodeNotifierBlock(block, ctx)
		if notifierDiags.HasErrors() {
			return taskFile, diags.Extend(notifierDiags)
		}

		if _, found := taskFile.Notifiers[notifier.Name]; found {
			return taskFile, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate notifier",
				Detail:   "Duplicate " + notifier.Name + " notifier definition found.",
				Subject:  &notifier.DeclRange,
				Context:  block.DefRange.Ptr(),
			})
		}

		taskFile.Notifiers[notifier.Name] = notifier
	}

	for _, block := range content.Blocks.OfType("task") {
		task, taskDiags := decodeTaskBlock(block, ctx)
		if taskDiags.HasErrors() {
			return taskFile, diags.Extend(taskDiags)
		}

		if _, found := taskFile.Tasks[task.Name]; found {
			return taskFile, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate task",
				Detail:   "Duplicate " + task.Name + " task definition found.",
				Subject:  &task.DeclRange,
				Context:  block.DefRange.Ptr(),
			})
		}

		taskFile.Tasks[task.Name] = task
	}

	return taskFile, diags
}
