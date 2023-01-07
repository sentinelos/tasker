// Package compressor implementation that will compressor that extension/type.
package compressor

import (
	"os"
)

// Compressor defines the interface that must be implemented to add support for compressor a type.
type Compressor interface {
	// Compress should compress src to dst. dir specifies whether dst
	// is a directory or single file. src is guaranteed to be a single file
	// that exists. dst is not guaranteed to exist already.
	Compress(src, dst string, dir bool) error

	// Decompress should decompress src to dst. dir specifies whether dst
	// is a directory or single file. src is guaranteed to be a single file
	// that exists. dst is not guaranteed to exist already.
	Decompress(src, dst string, dir bool, umask os.FileMode) error
}

// Compressors is the mapping of extension to the Compressor implementation that will compressor that extension/type.
var Compressors map[string]Compressor
