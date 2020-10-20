package cmd

import (
	"errors"
	"fmt"
	"kool-dev/kool/cmd/presets"
	"kool-dev/kool/cmd/shell"

	"github.com/spf13/cobra"
)

// KoolPresetFlags holds the flags for the preset command
type KoolPresetFlags struct {
	Override bool
}

// KoolPreset holds handlers and functions to implement the preset command logic
type KoolPreset struct {
	DefaultKoolService
	Flags        *KoolPresetFlags
	parser       presets.Parser
	promptSelect shell.PromptSelect
}

// ErrPresetFilesAlreadyExists error for existing presets files
var ErrPresetFilesAlreadyExists = errors.New("some preset files already exist")

func init() {
	var (
		preset    = NewKoolPreset()
		presetCmd = NewPresetCommand(preset)
	)

	rootCmd.AddCommand(presetCmd)
}

// NewKoolPreset creates a new handler for preset logic
func NewKoolPreset() *KoolPreset {
	return &KoolPreset{
		*newDefaultKoolService(),
		&KoolPresetFlags{false},
		&presets.DefaultParser{Presets: presets.GetAll()},
		shell.NewPromptSelect(),
	}
}

// Execute runs the preset logic with incoming arguments.
func (i *KoolPreset) Execute(args []string) (err error) {
	var fileError, preset, language string

	if len(args) == 0 {
		if language, err = i.promptSelect.Ask("What language do you want to use", i.parser.GetLanguages()); err != nil {
			return
		}

		if preset, err = i.promptSelect.Ask("What preset do you want to use", i.parser.GetPresets(language)); err != nil {
			return
		}
	} else {
		preset = args[0]
	}

	if !i.parser.Exists(preset) {
		err = fmt.Errorf("Unknown preset %s", preset)
		return
	}

	i.Println("Preset", preset, "is initializing!")

	if !i.Flags.Override {
		existingFiles := i.parser.LookUpFiles(preset)
		for _, fileName := range existingFiles {
			i.Warning("Preset file ", fileName, " already exists.")
		}

		if len(existingFiles) > 0 {
			err = ErrPresetFilesAlreadyExists
			return
		}
	}

	if fileError, err = i.parser.WriteFiles(preset); err != nil {
		err = fmt.Errorf("Failed to write preset file %s: %v", fileError, err)
		return
	}

	i.Success("Preset ", preset, " initialized!")
	return
}

// NewPresetCommand initializes new kool preset command
func NewPresetCommand(preset *KoolPreset) (presetCmd *cobra.Command) {
	presetCmd = &cobra.Command{
		Use:   "preset [PRESET]",
		Short: "Initialize kool preset in the current working directory. If no preset argument is specified you will be prompted to pick among the existing options.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			preset.SetWriter(cmd.OutOrStdout())

			if err := preset.Execute(args); err != nil {
				if err.Error() == ErrPresetFilesAlreadyExists.Error() {
					preset.Warning("Some preset files already exist. In case you wanna override them, use --override.")
					preset.Exit(2)
				} else if err.Error() == shell.ErrPromptSelectInterrupted.Error() {
					preset.Warning("Operation Cancelled")
					preset.Exit(0)
				} else {
					preset.Error(err)
					preset.Exit(1)
				}
			}
		},
	}

	presetCmd.Flags().BoolVarP(&preset.Flags.Override, "override", "", false, "Force replace local existing files with the preset files")
	return
}
