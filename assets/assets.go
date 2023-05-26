// Package assets contains non-go assets - sometimes exposed as an embedfs
package assets

import "embed"

// UIFs is the filesystem storing static html contents
//
//go:embed ui/*
var UIFs embed.FS
