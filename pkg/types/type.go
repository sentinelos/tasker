// Package types defines types.
package types

import (
	"github.com/sentinelos/actions/pkg/reporter"
)

// Actions is a struct that actions.yaml gets decoded into.
//
//nolint:govet
type Actions struct {
	Policies []*PolicyDeclaration `yaml:"policies"`
	reporter reporter.Reporter
}

// PolicyDeclaration allows a user to declare an arbitrary type along with a spec that will be decoded
// into the appropriate concrete type.
//
//nolint:govet
type PolicyDeclaration struct {
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"spec"`
}
