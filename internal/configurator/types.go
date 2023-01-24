package configurator

import (
	"github.com/hashicorp/hcl/v2"
)

type (
	ConfigFile struct {
		Filename string            `yaml:"filename"`
		Source   *hcl.File         `yaml:"-"`
		RunsOn   map[string]*RunOn `yaml:"runs_on" hcl:"runs_on,block"`
	}

	RunOn struct {
		Name         string            `yaml:"name" hcl:"name"`
		Image        string            `json:"image" hcl:"image"`
		Description  string            `yaml:"description" hcl:"description,optional"`
		Environments map[string]string `yaml:"environments" hcl:"environments,optional"`
		Cpu          uint8             `json:"cpu" hcl:"cpu,optional"`
		Memory       string            `json:"memory" hcl:"memory,optional"`
		Platform     string            `json:"platform" hcl:"platform,optional"`
		User         string            `json:"user" hcl:"user,optional"`
		Shell        string            `json:"shell" hcl:"shell"`
		Workdir      string            `json:"workdir" hcl:"workdir,optional"`
		DeclRange    hcl.Range         `yaml:"-"`
	}
)
