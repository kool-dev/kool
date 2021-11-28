package commands

import (
	_ "embed"
	"kool-dev/kool/core/automate"
	"kool-dev/kool/core/presets"

	"github.com/spf13/cobra"
)

// KoolRecipe holds handlers and functions to implement the preset command logic
type KoolRecipe struct {
	DefaultKoolService
}

func AddKoolRecipe(root *cobra.Command) {
	var (
		recipe    = NewKoolRecipe()
		recipeCmd = NewRecipeCommand(recipe)
	)

	root.AddCommand(recipeCmd)
}

// NewKoolRecipe creates a new handler for preset logic
func NewKoolRecipe() *KoolRecipe {
	return &KoolRecipe{
		*newDefaultKoolService(),
	}
}

// Execute runs the add logic with incoming arguments.
func (p *KoolRecipe) Execute(args []string) (err error) {
	var recipe string

	if len(args) == 0 {
		// no recipe; let's just print them all
		var metas []*automate.RecipeMetadata

		if metas, err = automate.GetRecipes(); err != nil {
			return
		}

		p.Shell().Warning("You need to provide the recipe name as argument.")
		p.Shell().Println("")
		p.Shell().Println("Available recipes:")

		for _, meta := range metas {
			if meta.Title != "" {
				p.Shell().Printf("  %s (%s)\n", meta.Slug, meta.Title)
			} else {
				p.Shell().Printf("  %s\n", meta.Slug)
			}
		}
		return
	}

	recipe = args[0]

	err = presets.NewParser().Add(recipe, p.Shell())

	return
}

// NewRecipeCommand initializes new kool add command
func NewRecipeCommand(recipe *KoolRecipe) (recipeCmd *cobra.Command) {
	recipeCmd = &cobra.Command{
		Use:                   "recipe [RECIPE]",
		Short:                 "Adds configuration for some recipe in the current work directory.",
		Long:                  `Run the defines steps for a recipe which can add/edit files the current project directory in order to add some new service or configuration.`,
		Args:                  cobra.MaximumNArgs(1),
		RunE:                  DefaultCommandRunFunction(recipe),
		DisableFlagsInUseLine: true,
	}

	return
}
