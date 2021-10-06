package main

import (
	"embed"
	"kool-dev/kool/core/automate"
	"kool-dev/kool/core/presets"
)

//go:embed presets/* templates/*
var source embed.FS

//go:embed recipes/*
var recipes embed.FS

func init() {
	presets.SetSource(source)
	automate.SetRecipesSource(recipes)
}
