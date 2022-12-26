package commands

import (
	"fmt"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/presets"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

// TODO: create flag for --no-preset so the command runs only the create portion of the preset config

// KoolCreate holds handlers and functions to implement the create command logic
type KoolCreate struct {
	DefaultKoolService
	parser presets.Parser
	env    environment.EnvStorage
}

func AddKoolCreate(root *cobra.Command) {
	var (
		createCmd = NewCreateCommand(NewKoolCreate())
	)

	root.AddCommand(createCmd)
}

// NewKoolCreate creates a new handler for create logic
func NewKoolCreate() *KoolCreate {
	return &KoolCreate{
		*newDefaultKoolService(),
		presets.NewParser(),
		environment.NewEnvStorage(),
	}
}

// Execute runs the create logic with incoming arguments.
func (c *KoolCreate) Execute(args []string) (err error) {
	var (
		preset          = args[0]
		createDirectory = args[1]
	)

	// sets env variable CREATE_DIRECTORY so preset can use it
	c.env.Set("CREATE_DIRECTORY", createDirectory)

	if !c.parser.Exists(preset) {
		err = fmt.Errorf("unknown preset %s", preset)
		return
	}

	c.Shell().Println("Creating new", preset, "project...")

	c.parser.PrepareExecutor(c.Shell())

	if err = c.parser.Create(preset); err != nil {
		return
	}

	c.Shell().Println("Initializing", preset, "preset...")

	if !path.IsAbs(createDirectory) {
		if createDirectory, err = filepath.Abs(createDirectory); err != nil {
			return
		}
	}

	if err = os.Chdir(createDirectory); err != nil {
		return
	}

	c.env.Set("PWD", createDirectory)

	if err = c.parser.Install(preset); err != nil {
		return
	}

	c.Shell().Success("Preset ", preset, " created successfully!")

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
