package executor

import (
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

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
	context.AddObjectVariable("executor", RunnerContext{
		Name:    "local",
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		User:    os.Getenv("USER"),
		Shell:   os.Getenv("SHELL"),
		Workdir: workdir,
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

func (c *Context) AddObjectVariable(varName string, v any) {
	out := make(map[string]string)

	structure := reflect.ValueOf(v)
	if structure.Kind() == reflect.Ptr {
		structure = structure.Elem()
	}

	if structure.Kind() != reflect.Struct {
		c.AddMapVariable(varName, out)
		return
	}

	ty := structure.Type()
	for i := 0; i < structure.NumField(); i++ {
		if tag := ty.Field(i).Tag.Get("ctx"); tag != "" && structure.Field(i).Kind() == reflect.String {
			out[tag] = structure.Field(i).String()
		}
	}

	c.AddMapVariable(varName, out)
}

func (c *Context) AddVariables(v map[string]cty.Value) {
	c.AddVariable("variable", cty.ObjectVal(v))
}

func (c *Context) AddRunner(v RunnerContext) {
	c.AddObjectVariable("executor", v)
}

type RunnerContext struct {
	Name    string `ctx:"name"`
	OS      string `ctx:"os"`
	Arch    string `ctx:"arch"`
	User    string `ctx:"user"`
	Shell   string `ctx:"shell"`
	Workdir string `ctx:"workdir"`
}

type GitContext struct {
	Event      string `json:"event" ctx:"event"`
	Actor      string `json:"actor" ctx:"actor"`
	Repository string `json:"repository" ctx:"repository"`
	Branch     string `json:"branch" ctx:"branch"`
	Reference  string `json:"reference" ctx:"reference"`
}

type TaskContext struct {
	Status     Status        `ctx:"status"`
	StartedAt  time.Duration `ctx:"started_at"`
	FinishedAt time.Duration `ctx:"finished_at"`
	Container  struct {
		ID      string          `ctx:"id"`
		Status  ContainerStatus `ctx:"status"`
		Network string          `ctx:"network"`
	} `ctx:"container"`
	Services map[string]struct {
		ID      string          `ctx:"id"`
		Status  ContainerStatus `ctx:"status"`
		Network string          `ctx:"network"`
	} `ctx:"services"`
}

type StepContext struct {
	Status     Status            `ctx:"status"`
	Outputs    map[string]string `ctx:"outputs"`
	StartedAt  time.Duration     `ctx:"started_at"`
	FinishedAt time.Duration     `ctx:"finished_at"`
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
