// Package ui is used to embed the frontend in the golang binary
package ui

import "embed"

//go:embed dist/*
var Assets embed.FS
