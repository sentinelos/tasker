package configurator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// LoadConfigFile is a wrapper around DecodeConfigFile that first reads the given filename from disk and parses.
func LoadConfigFile(filename string) (*ConfigFile, hcl.Diagnostics) {
	src, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Config file not found",
				Detail:   fmt.Sprintf("The Config file %s does not exist.", filename),
				Subject:  &hcl.Range{Filename: filename},
			}}
		}
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read Config file",
			Detail:   fmt.Sprintf("The Config file %s could not be read, %s.", filename, err),
			Subject:  &hcl.Range{Filename: filename},
		}}
	}

	if len(src) == 0 {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Empty Config file",
			Detail:   fmt.Sprintf("The Config file %s is empty.", filename),
			Subject:  &hcl.Range{Filename: filename},
		}}
	}

	var (
		file  *hcl.File
		diags hcl.Diagnostics
	)

	parser := hclparse.NewParser()
	switch suffix := strings.ToLower(filepath.Ext(filename)); suffix {
	case ".cf":
		file, diags = parser.ParseHCL(src, filename)
	case ".cf.json":
		file, diags = parser.ParseJSON(src, filename)
	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unsupported Config file format",
			Detail:   fmt.Sprintf("The Config file %s could not be read, unrecognized format suffix %s.", filename, suffix),
			Subject:  &hcl.Range{Filename: filename},
		})
	}

	if diags.HasErrors() {
		return &ConfigFile{
			Filename: filename,
			Source:   file,
		}, diags
	}

	return DecodeConfigFile(filename, file)
}

// DecodeConfigFile decodes and evaluates expressions in the given ConfigFile source code.
//
// The "filename" argument provides ConfigFile filename
// The "file" argument provides parsed ConfigFile
func DecodeConfigFile(filename string, file *hcl.File) (*ConfigFile, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	configFile := &ConfigFile{
		Filename: filename,
		Source:   file,
		RunsOn:   map[string]*RunOn{},
	}

	content, contentDiags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "run_on", LabelNames: []string{"name", "type"}},
		},
	})

	diags = diags.Extend(contentDiags)

	for _, block := range content.Blocks.OfType("run_on") {
		runOn, runOnDiags := decodeRunOnBlock(block)
		if runOnDiags.HasErrors() {
			return configFile, diags.Extend(runOnDiags)
		}

		if _, found := configFile.RunsOn[runOn.Name]; found {
			return configFile, diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate run_on",
				Detail:   "Duplicate " + runOn.Name + " run_on definition found.",
				Subject:  &runOn.DeclRange,
				Context:  block.DefRange.Ptr(),
			})
		}

		configFile.RunsOn[runOn.Name] = runOn
	}

	return configFile, diags
}
