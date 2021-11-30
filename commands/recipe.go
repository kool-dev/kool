package commands

import (
	_ "embed"
	"kool-dev/kool/core/automate"
	"kool-dev/kool/core/presets"
	"strings"

	"github.com/agnivade/levenshtein"
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

	recipe = args[0]

	err = presets.NewParser().Add(recipe, p.Shell())

	return
}

// NewRecipeCommand initializes new kool add command
func NewRecipeCommand(recipe *KoolRecipe) (recipeCmd *cobra.Command) {
	recipeCmd = &cobra.Command{
		Use:   "recipe [RECIPE]",
		Short: "Adds configuration for some recipe in the current work directory.",
		Long:  `Run the defines steps for a recipe which can add/edit files the current project directory in order to add some new service or configuration.`,
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var recipes []string
			metas, _ := automate.GetRecipes()
			for _, meta := range metas {
				if meta.Slug == args[0] {
					return nil, cobra.ShellCompDirectiveDefault
				}

				if strings.HasPrefix(meta.Slug, args[0]) || levenshtein.ComputeDistance(meta.Slug, args[0]) <= 2 {
					recipes = append(recipes, meta.Slug)
				}
			}

			return recipes, cobra.ShellCompDirectiveDefault
		},
		RunE:                  DefaultCommandRunFunction(recipe),
		DisableFlagsInUseLine: true,
	}

	return
}
