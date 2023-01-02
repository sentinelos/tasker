package taskfile

import (
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Status uint
type ContainerStatus uint

const (
	StatusPending Status = iota
	StatusSkipped
	StatusRunning
	StatusCancelled
	StatusSuccess
	StatusFailure
)

const (
	ContainerStatusUnknown ContainerStatus = iota
	ContainerStatusCreated
	ContainerStatusRunning
	ContainerStatusPausing
	ContainerStatusPaused
	ContainerStatusStopped
	ContainerStatusExited
)

var (
	StatusNames = map[Status]string{
		StatusPending:   "pending",
		StatusSkipped:   "skipped",
		StatusRunning:   "running",
		StatusCancelled: "cancelled",
		StatusSuccess:   "success",
		StatusFailure:   "failure",
	}

	ContainerStatusNames = map[ContainerStatus]string{
		ContainerStatusUnknown: "unknown",
		ContainerStatusCreated: "created",
		ContainerStatusRunning: "running",
		ContainerStatusPausing: "pausing",
		ContainerStatusPaused:  "paused",
		ContainerStatusStopped: "stopped",
		ContainerStatusExited:  "exited",
	}

	DefaultStatus          = StatusPending
	DefaultContainerStatus = ContainerStatusUnknown
)

type Context struct {
	Ctx *hcl.EvalContext
}

func NewContext() *Context {
	id := uuid.New().String()
	workdir, _ := os.Getwd()

	context := &Context{
		Ctx: &hcl.EvalContext{
			Variables: map[string]cty.Value{},
			Functions: map[string]function.Function{},
		},
	}

	context.AddStringVariable("run_id", id)
	context.AddObjectVariable("runner", RunnerContext{
		Name:    "local",
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		User:    os.Getenv("USER"),
		Shell:   os.Getenv("SHELL"),
		Workdir: workdir,
	}, map[string]cty.Type{
		"name":    cty.String,
		"os":      cty.String,
		"arch":    cty.String,
		"user":    cty.String,
		"shell":   cty.String,
		"workdir": cty.String,
	})

	return context
}

func (c *Context) AddVariable(varName string, value cty.Value) {
	c.Ctx.Variables[varName] = value
}

func (c *Context) AddStringVariable(varName string, v string) {
	value, _ := gocty.ToCtyValue(v, cty.String)
	c.AddVariable(varName, value)
}

func (c *Context) AddMapVariable(varName string, v map[string]string) {
	value, _ := gocty.ToCtyValue(v, cty.Map(cty.String))
	c.AddVariable(varName, value)
}

func (c *Context) AddListVariable(varName string, v []string) {
	value, _ := gocty.ToCtyValue(v, cty.List(cty.String))
	c.AddVariable(varName, value)
}

func (c *Context) AddObjectVariable(varName string, v any, ty map[string]cty.Type) {
	value, _ := gocty.ToCtyValue(v, cty.Object(ty))
	c.AddVariable(varName, value)
}

func (c *Context) AddVariables(v map[string]cty.Value) {
	c.AddVariable("variable", cty.ObjectVal(v))
}

type RunnerContext struct {
	Name    string `cty:"name"`
	OS      string `cty:"os"`
	Arch    string `cty:"arch"`
	User    string `cty:"user"`
	Shell   string `cty:"shell"`
	Workdir string `cty:"workdir"`
}

type GitContext struct {
	Event      string `json:"event" cty:"event"`
	Actor      string `json:"actor" cty:"actor"`
	Repository string `json:"repository" cty:"repository"`
	Branch     string `json:"branch" cty:"branch"`
	Reference  string `json:"reference" cty:"reference"`
}

type TaskContext struct {
	Status     Status        `cty:"status"`
	StartedAt  time.Duration `cty:"started_at"`
	FinishedAt time.Duration `cty:"finished_at"`
	Container  struct {
		ID      string          `cty:"id"`
		Status  ContainerStatus `cty:"status"`
		Network string          `cty:"network"`
	} `cty:"container"`
	Services map[string]struct {
		ID      string          `cty:"id"`
		Status  ContainerStatus `cty:"status"`
		Network string          `cty:"network"`
	} `cty:"services"`
}

type StepContext struct {
	Status     Status            `cty:"status"`
	Outputs    map[string]string `cty:"outputs"`
	StartedAt  time.Duration     `cty:"started_at"`
	FinishedAt time.Duration     `cty:"finished_at"`
}

// String status to string
func (s Status) String() string {
	return StatusNames[s]
}

// MarshalYAML status to yaml
func (s Status) MarshalYAML() ([]byte, error) {
	return []byte(s.String()), nil
}

// String container status to string
func (s ContainerStatus) String() string {
	return ContainerStatusNames[s]
}

// MarshalYAML container status to yaml
func (s ContainerStatus) MarshalYAML() ([]byte, error) {
	return []byte(s.String()), nil
}
