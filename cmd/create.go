package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/presets"
	"os"

	"github.com/spf13/cobra"
)

// KoolCreate holds handlers and functions to implement the preset command logic
type KoolCreate struct {
	DefaultKoolService
	parser        presets.Parser
	createCommand builder.Command
	KoolPreset
}

func init() {
	var (
		create    = NewKoolCreate()
		createCmd = NewCreateCommand(create)
	)

	rootCmd.AddCommand(createCmd)
}

// NewKoolCreate creates a new handler for create logic
func NewKoolCreate() *KoolCreate {
	return &KoolCreate{
		*newDefaultKoolService(),
		presets.NewParser(),
		&builder.DefaultCommand{},
		*NewKoolPreset(),
	}
}

// Execute runs the create logic with incoming arguments.
func (c *KoolCreate) Execute(originalArgs []string) (err error) {
	var (
		presetConfig *presets.PresetConfig
		createCmds   []string
		ok           bool
	)

	preset := originalArgs[0]
	dir := originalArgs[1]

	c.parser.LoadPresets(presets.GetAll())
	c.parser.LoadConfigs(presets.GetConfigs())

	if !c.parser.Exists(preset) {
		err = fmt.Errorf("Unknown preset %s", preset)
		return
	}

	if presetConfig, err = c.parser.GetConfig(preset); err != nil || presetConfig == nil {
		err = fmt.Errorf("error parsing preset config; err: %v", err)
		return
	}

	if createCmds, ok = presetConfig.Commands["create"]; !ok || len(createCmds) <= 0 {
		err = fmt.Errorf("No create commands were found for preset %s", preset)
		return
	}

	for _, createCmd := range createCmds {
		if err = c.createCommand.Parse(createCmd); err != nil {
			return
		}

		if err = c.Interactive(c.createCommand, dir); err != nil {
			return
		}
	}

	_ = os.Chdir(dir)

	err = c.KoolPreset.Execute([]string{preset})

	return
}

// NewCreateCommand initializes new kool create command
func NewCreateCommand(create *KoolCreate) (createCmd *cobra.Command) {
	createCmd = &cobra.Command{
		Use:   "create [preset] [project]",
		Short: "Create a new project using preset",
		Args:  cobra.ExactArgs(2),
		Run:   DefaultCommandRunFunction(create),
	}

	return
}
