package commands

import (
	"fmt"
	"kool-dev/kool/core/automate"
	"kool-dev/kool/core/presets"
	"kool-dev/kool/core/shell"
	"sort"

	"github.com/spf13/cobra"
)

// KoolPreset holds handlers and functions to implement the preset command logic
type KoolPreset struct {
	DefaultKoolService
	presetsParser presets.Parser
	promptSelect  shell.PromptSelect

	Flags *KoolAnswersFlags
}

// KoolAnswersFlags holds the flags for answering the preset questions
// beforehand, so they are not prompted interactively
type KoolAnswersFlags struct {
	Answers       []string
	NoInteraction bool
}

// AddAnswersFlags registers the flags for answering preset questions
func AddAnswersFlags(cmd *cobra.Command, flags *KoolAnswersFlags) {
	cmd.Flags().StringArrayVarP(&flags.Answers, "answer", "a", []string{}, "Answer to a preset question as REF=value (repeatable; comma separate multiple choices)")
	cmd.Flags().BoolVarP(&flags.NoInteraction, "no-interaction", "n", false, "Fail on any preset question without an answer instead of prompting")
}

// parse parses the provided answers flags
func (f *KoolAnswersFlags) parse() (answers automate.Answers, noInteraction bool, err error) {
	if f == nil {
		return
	}

	answers, err = automate.ParseAnswers(f.Answers)
	noInteraction = f.NoInteraction
	return
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
		shell.NewPromptSelect(),
		&KoolAnswersFlags{},
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

	answers, noInteraction, err := p.Flags.parse()
	if err != nil {
		return
	}

	p.Shell().Println("Preset", preset, "is initializing!")

	p.presetsParser.PrepareExecutor(p.Shell(), answers, noInteraction)

	if err = p.presetsParser.Install(preset); err != nil {
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
an interactive wizard will present the available options.

Questions asked by the preset can be answered beforehand with --answer REF=value,
where REF is the question ref (or the question text when it has no ref). Options
can be given by value or label, and multiple choices are comma separated. Without
a terminal, or with --no-interaction, a question without an answer fails.`,
		Args:                  cobra.MaximumNArgs(1),
		RunE:                  DefaultCommandRunFunction(preset),
		DisableFlagsInUseLine: true,
	}

	AddAnswersFlags(presetCmd, preset.Flags)

	return
}

func (p *KoolPreset) getPreset(args []string) (pickedPreset string, err error) {
	if len(args) == 1 {
		pickedPreset = args[0]
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

	var availablePresets = p.presetsParser.GetPresets(tag)
	var presets []string

	for _, name := range availablePresets {
		presets = append(presets, name)
	}

	sort.Strings(presets)

	if pickedPreset, err = p.promptSelect.Ask("What preset do you want to use", presets); err != nil {
		return
	}

	for preset, name := range availablePresets {
		if name == pickedPreset {
			pickedPreset = preset
		}
	}

	return
}
