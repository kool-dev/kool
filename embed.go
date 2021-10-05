package main

import (
	"embed"
	"kool-dev/kool/core/presets"
)

//go:embed presets/* templates/*
var source embed.FS

func init() {
	presets.SetSource(source)
}
