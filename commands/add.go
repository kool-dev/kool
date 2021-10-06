package commands

import (
	_ "embed"
	"kool-dev/kool/core/presets"

	"github.com/spf13/cobra"
)

// KoolAdd holds handlers and functions to implement the preset command logic
type KoolAdd struct {
	DefaultKoolService
}

func AddKoolAdd(root *cobra.Command) {
	var (
		add    = NewKoolAdd()
		addCmd = NewAddCommand(add)
	)

	root.AddCommand(addCmd)
}

// NewKoolAdd creates a new handler for preset logic
func NewKoolAdd() *KoolAdd {
	return &KoolAdd{
		*newDefaultKoolService(),
	}
}

// Execute runs the add logic with incoming arguments.
func (p *KoolAdd) Execute(args []string) (err error) {
	var recipe = args[0]

	err = presets.NewParser().Add(recipe, p.Shell())

	return
}

// NewAddCommand initializes new kool add command
func NewAddCommand(add *KoolAdd) (addCmd *cobra.Command) {
	addCmd = &cobra.Command{
		Use:                   "add [RECIPE]",
		Short:                 "Adds configuration for some recipe in the current work directory.",
		Long:                  `Run the defines steps for a recipe which can add/edit files the current project directory in order to add some new service or configuration.`,
		Args:                  cobra.ExactArgs(1),
		RunE:                  DefaultCommandRunFunction(add),
		DisableFlagsInUseLine: true,
	}

	return
}
