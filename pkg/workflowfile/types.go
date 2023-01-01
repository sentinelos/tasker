package workflowfile

import (
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// Workflow is the structure of the files in .workflows
type Workflow struct {
	Filename    string             `yaml:"filename"`
	Name        string             `yaml:"name" hcl:"name"`
	Description string             `yaml:"description" hcl:"description,optional"`
	Notifies    map[string]*Notify `yaml:"notifies" hcl:"notify,block"`
	Jobs        map[string]*Job    `yaml:"jobs" hcl:"job,block"`
}

type Job struct {
	Name        string                `yaml:"name" hcl:"name"`
	Description string                `yaml:"description" hcl:"description,optional"`
	Trigger     map[string]*Trigger   `yaml:"trigger" hcl:"trigger,block"`
	DependsOn   []string              `yaml:"depends_on" hcl:"depends_on,optional"`
	Uses        *Use                  `yaml:"uses" hcl:"use,block"`
	Container   *Container            `yaml:"container" hcl:"container,block"`
	Services    map[string]*Container `yaml:"services" hcl:"service,block"`
	Steps       map[string]*Step      `yaml:"steps" hcl:"step,block"`
	Timeout     time.Duration         `yaml:"timeout" hcl:"timeout,optional"`
	DeclRange   hcl.Range             `yaml:"-"`
}

type Step struct {
	Name        string        `yaml:"name" hcl:"name"`
	Description string        `yaml:"description" hcl:"description,optional"`
	If          bool          `yaml:"if" hcl:"if,optional"`
	OnFailure   string        `yaml:"on_failure" hcl:"on_failure,optional"`
	Uses        *Use          `yaml:"uses" hcl:"use,block"`
	Shell       string        `yaml:"shell" hcl:"shell,optional"`
	Workdir     string        `yaml:"workdir" hcl:"workdir,optional"`
	Command     string        `yaml:"command" hcl:"command"`
	Timeout     time.Duration `yaml:"timeout" hcl:"timeout,optional"`
	DeclRange   hcl.Range     `yaml:"-"`
}

type Notify struct {
	Name        string             `yaml:"name" hcl:"name"`
	Description string             `yaml:"description" hcl:"description,optional"`
	Outputs     map[string]*Output `yaml:"outputs" hcl:"output,block"`
	DeclRange   hcl.Range          `yaml:"-"`
}

type Trigger struct {
	Name        string    `yaml:"name" hcl:"name"`
	Description string    `yaml:"description" hcl:"description,optional"`
	Conditions  []string  `yaml:"conditions" hcl:"conditions"`
	DeclRange   hcl.Range `yaml:"-"`
}

type Use struct {
	Name        string            `yaml:"name" hcl:"name"`
	Description string            `yaml:"description" hcl:"description,optional"`
	Inputs      map[string]*Input `yaml:"inputs" hcl:"input,block"`
	DeclRange   hcl.Range         `yaml:"-"`
}

type Variable struct {
	Name        string    `yaml:"name" hcl:"name"`
	Description string    `yaml:"description" hcl:"description,optional"`
	Value       cty.Value `yaml:"value" hcl:"value"`
	Sensitive   bool      `yaml:"sensitive" hcl:"sensitive,optional"`
	DeclRange   hcl.Range `yaml:"-"`
}

type Input struct {
	Name        string    `yaml:"name" hcl:"name"`
	Description string    `yaml:"description" hcl:"description,optional"`
	Value       string    `yaml:"value" hcl:"value"`
	Sensitive   bool      `yaml:"sensitive" hcl:"sensitive,optional"`
	DeclRange   hcl.Range `yaml:"-"`
}

type Output struct {
	Name        string    `yaml:"name" hcl:"name"`
	Description string    `yaml:"description" hcl:"description,optional"`
	Value       string    `yaml:"value" hcl:"value"`
	Sensitive   bool      `yaml:"sensitive" hcl:"sensitive,optional"`
	DeclRange   hcl.Range `yaml:"-"`
}

// Container is the specification of the container to use for the job
type Container struct {
	Image        string            `yaml:"image" hcl:"image"`
	Description  string            `yaml:"description" hcl:"description,optional"`
	Environments map[string]string `yaml:"environments" hcl:"environments,optional"`
	Volumes      []string          `yaml:"volumes" hcl:"volumes,optional"`
	Flags        string            `yaml:"flags" hcl:"flags,optional"`
	Command      string            `yaml:"command" hcl:"command,optional"`
	Args         string            `yaml:"args" hcl:"args,optional"`
	Credential   *Credential       `yaml:"credential" hcl:"credential,block"`
	DeclRange    hcl.Range         `yaml:"-"`
}

type Credential struct {
	Username  string    `yaml:"username" hcl:"username"`
	Password  string    `yaml:"password" hcl:"password"`
	DeclRange hcl.Range `yaml:"-"`
}
