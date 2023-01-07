// Package taskfile defines file specification.
package taskfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

const (
	// BadIdentifierDetail A consistent detail message for all "not a valid identifier" diagnostics.
	BadIdentifierDetail = "A name must start with a letter or underscore and may contain only letters, digits, underscores, and dashes."

	// DefaultTimeout is used if there is no timeout given
	DefaultTimeout = 10 * time.Minute
)

// LoadTaskfile is a wrapper around DecodeTaskfile that first reads the given filename from disk and parses.
func LoadTaskfile(filename string, context *Context) (*Taskfile, hcl.Diagnostics) {
	src, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Taskfile not found",
				Detail:   fmt.Sprintf("The Taskfile %s does not exist.", filename),
				Subject:  &hcl.Range{Filename: filename},
			}}
		}
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read Taskfile",
			Detail:   fmt.Sprintf("The Taskfile %s could not be read, %s.", filename, err),
			Subject:  &hcl.Range{Filename: filename},
		}}
	}

	if len(src) == 0 {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Empty Taskfile",
			Detail:   fmt.Sprintf("The Taskfile %s is empty.", filename),
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
			Summary:  "Unsupported Taskfile format",
			Detail:   fmt.Sprintf("The Taskfile %s could not be read, unrecognized format suffix %s.", filename, suffix),
			Subject:  &hcl.Range{Filename: filename},
		})
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return DecodeTaskfile(filename, file, context)
}

// DecodeTaskfile decodes and evaluates expressions in the given Taskfile source code.
//
// The "filename" argument provides Taskfile filename
// The "file" argument provides parsed Taskfile
// The "context" argument provides variables and functions for use during expression evaluation.
func DecodeTaskfile(filename string, file *hcl.File, context *Context) (*Taskfile, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	taskfile := &Taskfile{
		Filename: filename,
		Notifies: map[string]*Notify{},
		Tasks:    map[string]*Task{},
	}

	content, contentDiags := file.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "name", Required: true},
			{Name: "description"},
			{Name: "runs_on", Required: true},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "notify", LabelNames: []string{"name"}},
			{Type: "task", LabelNames: []string{"name"}},
		},
	})

	diags = diags.Extend(contentDiags)

	if attr, exists := content.Attributes["name"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &taskfile.Name))
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &taskfile.Description))
	}

	if attr, exists := content.Attributes["runs_on"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &taskfile.RunsOn))
	}

	variables := map[string]cty.Value{}
	for _, block := range content.Blocks.OfType("variable") {
		variable, varDiags := decodeVariableBlock(block, nil)
		if varDiags.HasErrors() {
			return taskfile, diags.Extend(varDiags)
		}

		if _, found := variables[variable.Name]; found {
			return taskfile, diags.Append(&hcl.Diagnostic{
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
	for _, block := range content.Blocks.OfType("notify") {
		notify, notifyDiags := decodeNotifyBlock(block, ctx)
		if notifyDiags.HasErrors() {
			return taskfile, diags.Extend(notifyDiags)
		}

		if _, found := taskfile.Notifies[notify.Name]; found {
			return taskfile, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate notify",
				Detail:   "Duplicate " + notify.Name + " notify definition found.",
				Subject:  &notify.DeclRange,
				Context:  block.DefRange.Ptr(),
			})
		}

		taskfile.Notifies[notify.Name] = notify
	}

	for _, block := range content.Blocks.OfType("task") {
		task, taskDiags := decodeTaskBlock(block, ctx)
		if taskDiags.HasErrors() {
			return taskfile, diags.Extend(taskDiags)
		}

		if _, found := taskfile.Tasks[task.Name]; found {
			return taskfile, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate task",
				Detail:   "Duplicate " + task.Name + " task definition found.",
				Subject:  &task.DeclRange,
				Context:  block.DefRange.Ptr(),
			})
		}

		taskfile.Tasks[task.Name] = task
	}

	return taskfile, diags
}
