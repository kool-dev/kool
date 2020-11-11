package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/presets"
	"strings"

	"github.com/spf13/cobra"
)

// KoolCreate holds handlers and functions to implement the preset command logic
type KoolCreate struct {
	DefaultKoolService
	parser     presets.Parser
	koolDocker builder.Command
}

func init() {
	var (
		create    = NewKoolCreate()
		createCmd = NewCreateCommand(create)
	)

	rootCmd.AddCommand(createCmd)
}

// NewKoolCreate creates a new handler for exec logic
func NewKoolCreate() *KoolCreate {
	return &KoolCreate{
		*newDefaultKoolService(),
		&presets.DefaultParser{Presets: presets.GetAll()},
		builder.NewCommand("kool", "docker"),
	}
}

// Execute runs the exec logic with incoming arguments.
func (c *KoolCreate) Execute(args []string) (err error) {
	if !c.IsTerminal() {
		c.koolDocker.AppendArgs("-T")
	}

	createCmd, err := c.parser.GetCreateCommand(args[0])

	if err != nil {
		return
	}

	c.koolDocker.AppendArgs(strings.Fields(createCmd)...)

	err = c.koolDocker.Interactive(args[1:]...)

	return
}

// NewCreateCommand initializes new kool create command
func NewCreateCommand(create *KoolCreate) (createCmd *cobra.Command) {
	createCmd = &cobra.Command{
		Use:   "create [preset] [project]",
		Short: "Create a new project using preset",
		Args:  cobra.MinimumNArgs(2),
		Run:   DefaultCommandRunFunction(create),
	}

	return
}
