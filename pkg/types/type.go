// Package types defines types.
package types

import (
	"github.com/sentinelos/ensurer/pkg/reporter"
)

// Ensurer is a struct that ensurer.yaml gets decoded into.
//
//nolint:govet
type Ensurer struct {
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
