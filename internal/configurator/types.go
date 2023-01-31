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
		Type         string            `yaml:"type" hcl:"type"`
		Description  string            `yaml:"description" hcl:"description,optional"`
		Environments map[string]string `yaml:"environments" hcl:"environments,optional"`
		Image        string            `json:"image" hcl:"image"`
		Labels       map[string]string `yaml:"labels" hcl:"labels,optional"`
		Cpu          int64             `json:"cpu" hcl:"cpu,optional"`
		Memory       int64             `json:"memory" hcl:"memory,optional"`
		Platform     string            `json:"platform" hcl:"platform,optional"`
		User         string            `json:"user" hcl:"user,optional"`
		Shell        string            `json:"shell" hcl:"shell"`
		Workdir      string            `json:"workdir" hcl:"workdir,optional"`
		DeclRange    hcl.Range         `yaml:"-"`
	}
)
