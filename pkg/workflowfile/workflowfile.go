// Package workflowfile defines file specification.
package workflowfile

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

// LoadWorkflow is a wrapper around DecodeWorkflow that first reads the given filename
// from disk. See the DecodeWorkflow documentation for more information.
func LoadWorkflow(filename string, context *Context) (*Workflow, hcl.Diagnostics) {
	src, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Workflow file not found",
				Detail:   fmt.Sprintf("The workflow file does not exist."),
				Subject:  &hcl.Range{Filename: filename},
			}}
		}
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read workflow",
			Detail:   fmt.Sprintf("The workflow file could not be read: %s.", err),
			Subject:  &hcl.Range{Filename: filename},
		}}
	}

	if len(src) == 0 {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Empty workflow file",
			Detail:   fmt.Sprintf("Empty workflow file."),
			Subject:  &hcl.Range{Filename: filename},
		}}
	}

	return DecodeWorkflow(filename, src, context)
}

func DecodeWorkflow(filename string, src []byte, context *Context) (*Workflow, hcl.Diagnostics) {
	var (
		file  *hcl.File
		diags hcl.Diagnostics
	)

	parser := hclparse.NewParser()
	switch suffix := strings.ToLower(filepath.Ext(filename)); suffix {
	case ".wf":
		file, diags = parser.ParseHCL(src, filename)
	case ".wf.json":
		file, diags = parser.ParseJSON(src, filename)
	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported workflow file format",
			Detail:   fmt.Sprintf("Cannot read from %s: unrecognized file format suffix %q.", filename, suffix),
		})
		return nil, diags
	}

	workflow := &Workflow{
		Filename: filename,
		Notifies: map[string]*Notify{},
		Jobs:     map[string]*Job{},
	}

	content, contentDiags := file.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "name", Required: true},
			{Name: "description"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "notify", LabelNames: []string{"name"}},
			{Type: "job", LabelNames: []string{"name"}},
		},
	})

	diags = diags.Extend(contentDiags)

	variables := map[string]cty.Value{}
	for _, block := range content.Blocks.OfType("variable") {
		variable, varDiags := decodeVariableBlock(block, nil)
		if varDiags.HasErrors() {
			return nil, diags.Extend(varDiags)
		}

		if _, found := variables[variable.Name]; found {
			return nil, diags.Append(&hcl.Diagnostic{
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

	if attr, exists := content.Attributes["name"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &workflow.Name))
	}

	if attr, exists := content.Attributes["description"]; exists {
		diags = diags.Extend(gohcl.DecodeExpression(attr.Expr, nil, &workflow.Description))
	}

	ctx := context.Ctx
	for _, block := range content.Blocks.OfType("notify") {
		notify, notifyDiags := decodeNotifyBlock(block, ctx)
		if notifyDiags.HasErrors() {
			return nil, diags.Extend(notifyDiags)
		}

		if _, found := workflow.Notifies[notify.Name]; found {
			return nil, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate notify",
				Detail:   "Duplicate " + notify.Name + " notify definition found.",
				Subject:  &notify.DeclRange,
				Context:  block.DefRange.Ptr(),
			})
		}

		workflow.Notifies[notify.Name] = notify
	}

	for _, block := range content.Blocks.OfType("job") {
		job, jobDiags := decodeJobBlock(block, ctx)
		if jobDiags.HasErrors() {
			return nil, diags.Extend(jobDiags)
		}

		if _, found := workflow.Jobs[job.Name]; found {
			return nil, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate job",
				Detail:   "Duplicate " + job.Name + " job definition found.",
				Subject:  &job.DeclRange,
				Context:  block.DefRange.Ptr(),
			})
		}

		workflow.Jobs[job.Name] = job
	}

	return workflow, diags
}
