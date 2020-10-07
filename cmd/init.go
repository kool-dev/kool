package cmd

import (
	"errors"
	"fmt"
	"kool-dev/kool/cmd/presets"

	"github.com/spf13/cobra"
)

// KoolInitFlags holds the flags for the init command
type KoolInitFlags struct {
	Override bool
}

// KoolInit holds handlers and functions to implement the init command logic
type KoolInit struct {
	DefaultKoolService
	Flags  *KoolInitFlags
	parser presets.Parser
}

// ErrPresetFilesAlreadyExists error for existing presets files
var ErrPresetFilesAlreadyExists = errors.New("some preset files already exist")

func init() {
	var (
		init    = NewKoolInit()
		initCmd = NewInitCommand(init)
	)

	rootCmd.AddCommand(initCmd)
}

// NewKoolInit creates a new handler for init logic
func NewKoolInit() *KoolInit {
	return &KoolInit{
		*newDefaultKoolService(),
		&KoolInitFlags{false},
		&presets.DefaultParser{Presets: presets.GetAll()},
	}
}

// Execute runs the init logic with incoming arguments.
func (i *KoolInit) Execute(args []string) (err error) {
	var fileError, preset string

	preset = args[0]

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

	return
}

// NewInitCommand initializes new kool init command
func NewInitCommand(init *KoolInit) (initCmd *cobra.Command) {
	initCmd = &cobra.Command{
		Use:   "init [PRESET]",
		Short: "Initialize kool preset in the current working directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			init.SetWriter(cmd.OutOrStdout())

			if err := init.Execute(args); err != nil {
				if err.Error() == ErrPresetFilesAlreadyExists.Error() {
					init.Warning("Some preset files already exist. In case you wanna override them, use --override.")
					init.Exit(2)
				} else {
					init.Error(err)
					init.Exit(1)
				}
			} else {
				init.Success("Preset ", args[0], " initialized!")
			}
		},
	}

	initCmd.Flags().BoolVarP(&init.Flags.Override, "override", "", false, "Force replace local existing files with the preset files")
	return
}
