package commands

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/presets"

	"github.com/spf13/cobra"
)

// KoolCreate holds handlers and functions to implement the preset command logic
type KoolCreate struct {
	DefaultKoolService
	parser        presets.Parser
	env           environment.EnvStorage
	createCommand builder.Command
	KoolPreset
}

func AddKoolCreate(root *cobra.Command) {
	var (
		create    = NewKoolCreate()
		createCmd = NewCreateCommand(create)
	)

	root.AddCommand(createCmd)
}

// NewKoolCreate creates a new handler for create logic
func NewKoolCreate() *KoolCreate {
	return &KoolCreate{
		*newDefaultKoolService(),
		presets.NewParser(),
		environment.NewEnvStorage(),
		&builder.DefaultCommand{},
		*NewKoolPreset(),
	}
}

// Execute runs the create logic with incoming arguments.
func (c *KoolCreate) Execute(args []string) (err error) {
	var (
		preset          = args[0]
		createDirectory = args[1]
	)

	// sets env variable CREATE_DIRECTORY that aims to tell
	c.env.Set("CREATE_DIRECTORY", createDirectory)

	if !c.parser.Exists(preset) {
		err = fmt.Errorf("unknown preset %s", preset)
		return
	}

	// TODO: implement parser create run

	// for _, createCmd := range createCmds {
	// 	if err = c.createCommand.Parse(createCmd); err != nil {
	// 		return
	// 	}

	// 	if err = c.Shell().Interactive(c.createCommand); err != nil {
	// 		return
	// 	}
	// }

	// _ = os.Chdir(createDirectory)

	err = c.KoolPreset.Execute([]string{preset})

	return
}

// NewCreateCommand initializes new kool create command
func NewCreateCommand(create *KoolCreate) (createCmd *cobra.Command) {
	createCmd = &cobra.Command{
		Use:   "create PRESET FOLDER",
		Short: "Create a new project using a preset",
		Long:  "Create a new project using the specified PRESET in a directory named FOLDER.",
		Args:  cobra.ExactArgs(2),
		RunE:  DefaultCommandRunFunction(create),

		DisableFlagsInUseLine: true,
	}

	return
}
