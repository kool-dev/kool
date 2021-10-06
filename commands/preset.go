package commands

import (
	"fmt"
	"kool-dev/kool/core/parser"
	"kool-dev/kool/core/presets"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/core/templates"
	"kool-dev/kool/services/compose"

	"github.com/spf13/cobra"
)

// KoolPreset holds handlers and functions to implement the preset command logic
type KoolPreset struct {
	DefaultKoolService
	presetsParser  presets.Parser
	composeParser  compose.Parser
	templateParser templates.Parser
	koolYamlParser parser.KoolYamlParser
	promptSelect   shell.PromptSelect
}

func AddKoolPreset(root *cobra.Command) {
	var (
		preset    = NewKoolPreset()
		presetCmd = NewPresetCommand(preset)
	)

	root.AddCommand(presetCmd)
}

// NewKoolPreset creates a new handler for preset logic
func NewKoolPreset() *KoolPreset {
	return &KoolPreset{
		*newDefaultKoolService(),
		presets.NewParser(),
		compose.NewParser(),
		templates.NewParser(),
		&parser.KoolYaml{},
		shell.NewPromptSelect(),
	}
}

// Execute runs the preset logic with incoming arguments.
func (p *KoolPreset) Execute(args []string) (err error) {
	var preset string

	if preset, err = p.getPreset(args); err != nil {
		return
	}

	if !p.presetsParser.Exists(preset) {
		err = fmt.Errorf("unknown preset %s", preset)
		return
	}

	p.Shell().Println("Preset", preset, "is initializing!")

	if err = p.presetsParser.Install(preset, p.Shell()); err != nil {
		return
	}

	p.Shell().Success("Preset ", preset, " initialized!")
	return
}

// NewPresetCommand initializes new kool preset command
func NewPresetCommand(preset *KoolPreset) (presetCmd *cobra.Command) {
	presetCmd = &cobra.Command{
		Use:   "preset [PRESET]",
		Short: "Install configuration files customized for Kool in the current directory",
		Long: `Initialize a project using the specified [PRESET] by installing configuration
files customized for Kool in the current working directory. If no [PRESET] is provided,
an interactive wizard will present the available options.`,
		Args:                  cobra.MaximumNArgs(1),
		RunE:                  DefaultCommandRunFunction(preset),
		DisableFlagsInUseLine: true,
	}

	return
}

func (p *KoolPreset) getPreset(args []string) (preset string, err error) {
	if len(args) == 1 {
		preset = args[0]
		return
	}

	if !p.Shell().IsTerminal() {
		err = fmt.Errorf("please specify a preset as argument (non-TTY env)")
		return
	}

	var tag string
	if tag, err = p.promptSelect.Ask("Pick the preset category you are looking for", p.presetsParser.GetTags()); err != nil {
		return
	}

	preset, err = p.promptSelect.Ask("What preset do you want to use", p.presetsParser.GetPresets(tag))
	return
}
