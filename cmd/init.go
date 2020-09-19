package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/shell"
	"os"

	"github.com/spf13/cobra"
)

// InitFlags holds the flags for the init command
type InitFlags struct {
	Override bool
}

var initCmd = &cobra.Command{
	Use:   "init [PRESET]",
	Short: "Initialize kool preset in the current working directory",
	Args:  cobra.ExactArgs(1),
	Run:   runInit,
}

var initFlags = &InitFlags{false}

var initCmdOutputWriter shell.OutputWriter = shell.NewOutputWriter()

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&initFlags.Override, "override", "", false, "Force replace local existing files with the preset files")
}

func runInit(cmd *cobra.Command, args []string) {
	var (
		presetFiles                   map[string]string
		exists, hasExistingFile       bool
		preset, fileName, fileContent string
		err                           error
		file                          *os.File
		wrote                         int
	)

	initCmdOutputWriter.SetWriter(cmd.OutOrStdout())
	preset = args[0]

	if presetFiles, exists = presets[preset]; !exists {
		initCmdOutputWriter.Error("Unknown preset", preset)
		os.Exit(2)
	}

	fmt.Println("Preset", preset, "is initializing!")

	for fileName = range presetFiles {
		if !initFlags.Override {
			if _, err = os.Stat(fileName); !os.IsNotExist(err) {
				initCmdOutputWriter.Warning("  Preset file", fileName, "already exists.")
				hasExistingFile = true
			}
		}
	}

	if hasExistingFile {
		initCmdOutputWriter.Warning("Some preset files already exist. In case you wanna override them, use --override.")
		os.Exit(2)
	}

	for fileName, fileContent = range presetFiles {
		file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

		if err != nil {
			initCmdOutputWriter.Error("  Failed to create preset file", fileName, "due to error:", err)
			os.Exit(2)
		}

		if wrote, err = file.Write([]byte(fileContent)); err != nil {
			initCmdOutputWriter.Error("  Failed to write preset file", fileName, "due to error:", err)
			os.Exit(2)
		}

		if len([]byte(fileContent)) != wrote {
			initCmdOutputWriter.Error("  Failed to write preset file", fileName, " - failed to write all bytes:", wrote)
			os.Exit(2)
		}

		if err = file.Sync(); err != nil {
			initCmdOutputWriter.Error("  Failed to sync preset file", fileName, "due to error:", err)
			os.Exit(2)
		}

		file.Close()

		initCmdOutputWriter.Success("  Preset file", fileName, "created.")
	}

	initCmdOutputWriter.Success("Preset ", preset, " initialized!")
}
