package commands

import (
	"fmt"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/presets"
	"kool-dev/kool/core/shell"
	"os"
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
		createDirectory, preset string
	)

	if len(args) == 2 {
		preset = args[0]
		createDirectory = args[1]
	} else if len(args) == 1 {
		err = fmt.Errorf("bad number of arguments - either specify both preset and directory or none")
		return
	} else {
		if preset, err = NewKoolPreset().getPreset(args); err != nil {
			return
		}

		for {
			if createDirectory, err = shell.NewPromptInput().Input("New folder name:", fmt.Sprintf("my-kool-%s-project", preset)); err != nil {
				return
			}

			if createDirectory == "" {
				c.Shell().Error(fmt.Errorf("please enter a valid folder name"))
				continue
			} else if _, err = os.Stat(createDirectory); !os.IsNotExist(err) {
				c.Shell().Error(fmt.Errorf("folder %s already exists", createDirectory))
				continue
			} else {
				if err = os.MkdirAll(filepath.Join(os.TempDir(), createDirectory), 0755); err != nil {
					c.Shell().Error(fmt.Errorf("please enter a valid folder name"))
					continue
				} else {
					// ok we created, let's just have it removed if we fail
					defer func() {
						_ = os.RemoveAll(filepath.Join(os.TempDir(), createDirectory))
					}()
				}
			}

			// if no error, we got our directory
			break
		}
	}

	if !c.parser.Exists(preset) {
		err = fmt.Errorf("unknown preset %s", preset)
		return
	}

	var projectDir string
	if projectDir, err = c.prepareCreateTarget(createDirectory); err != nil {
		return
	}

	c.Shell().Println("Creating new", preset, "project...")

	c.parser.PrepareExecutor(c.Shell())

	if err = c.parser.Create(preset); err != nil {
		return
	}

	c.Shell().Println("Initializing", preset, "preset...")

	if _, err = os.Stat(projectDir); os.IsNotExist(err) {
		err = fmt.Errorf("preset did not create folder %s", projectDir)
		return
	}

	if err = os.Chdir(projectDir); err != nil {
		return
	}

	c.env.Set("PWD", projectDir)

	if err = c.parser.Install(preset); err != nil {
		return
	}

	c.Shell().Success("Preset ", preset, " created successfully!")

	return
}

// prepareCreateTarget makes CREATE_DIRECTORY a folder name relative to the
// current working directory so `kool docker` (which mounts cwd at /app) writes
// the new project onto the host. Absolute paths like /tmp/my-app otherwise
// land in the container's own /tmp and vanish when the container exits.
func (c *KoolCreate) prepareCreateTarget(createDirectory string) (projectDir string, err error) {
	if projectDir, err = filepath.Abs(createDirectory); err != nil {
		return
	}
	projectDir = filepath.Clean(projectDir)

	parent := filepath.Dir(projectDir)
	base := filepath.Base(projectDir)

	if err = os.MkdirAll(parent, 0755); err != nil {
		return
	}
	if err = os.Chdir(parent); err != nil {
		return
	}

	c.env.Set("CREATE_DIRECTORY", base)
	c.env.Set("PWD", parent)
	return
}

// NewCreateCommand initializes new kool create command
func NewCreateCommand(create *KoolCreate) (createCmd *cobra.Command) {
	createCmd = &cobra.Command{
		Use:   "create PRESET FOLDER",
		Short: "Create a new project using a preset",
		Long:  "Create a new project using the specified PRESET in a directory named FOLDER.",
		Args:  cobra.MaximumNArgs(2),
		RunE:  DefaultCommandRunFunction(create),

		DisableFlagsInUseLine: true,
	}

	return
}
