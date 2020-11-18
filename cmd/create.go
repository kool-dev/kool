package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/presets"
	"strings"

	"github.com/spf13/cobra"
)

// KoolCreate holds handlers and functions to implement the preset command logic
type KoolCreate struct {
	DefaultKoolService
	parser presets.Parser
	KoolDocker
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
		*NewKoolDocker(),
	}
}

// Execute runs the exec logic with incoming arguments.
func (c *KoolCreate) Execute(originalArgs []string) (err error) {
	preset := originalArgs[0]

	if !c.parser.Exists(preset) {
		err = fmt.Errorf("Unknown preset %s", preset)
		return
	}

	createCmd, err := c.parser.GetCreateCommand(preset)

	if err != nil {
		return
	}

	args := append(strings.Fields(createCmd), originalArgs[1:]...)

	c.KoolDocker.Execute(args)

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
