package automate

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v2"
)

type RecipeMetadata struct {
	Title string `yaml:"title"`
	Slug  string
}

var recipesSource embed.FS

func SetRecipesSource(src embed.FS) {
	recipesSource = src
}

func GetRecipes() (recipes []*RecipeMetadata, err error) {
	var (
		entries []fs.DirEntry
		raw     []byte
	)

	if entries, err = recipesSource.ReadDir("recipes"); err != nil {
		return
	}

	for _, e := range entries {
		m := &RecipeMetadata{Slug: strings.ReplaceAll(e.Name(), ".yml", "")}
		raw, err = recipesSource.ReadFile(fmt.Sprintf("recipes/%s", e.Name()))
		if err = yaml.Unmarshal(raw, m); err != nil {
			recipes = nil
			return
		}
		recipes = append(recipes, m)
	}
	return
}
