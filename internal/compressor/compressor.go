// Package compressor implementation that will compressor that extension/type.
package compressor

var (
	// Compressors is the mapping of extension to the Compressor implementation that will compressor that extension/type.
	Compressors map[string]Compressor
)
