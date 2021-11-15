package automate

import "embed"

var recipesSource embed.FS

func SetRecipesSource(src embed.FS) {
	recipesSource = src
}
