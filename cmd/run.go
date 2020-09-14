package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/parser"
	"kool-dev/kool/cmd/shell"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [SCRIPT]",
	Short: "Runs a custom command defined at kool.yaml in the working directory or in the kool folder of the user's home directory",
	Args:  cobra.MinimumNArgs(1),
	Run:   runRun,
}

func init() {
	rootCmd.AddCommand(runCmd)

	// after a non-flag arg, stop parsing flags
	runCmd.Flags().SetInterspersed(false)
}

func runRun(cmd *cobra.Command, args []string) {
	var (
		script   string
		commands []*builder.Command
		err      error
		kool     parser.Parser
	)

	kool = parser.NewParser()

	// look for kool.yml on current working directory
	_ = kool.AddLookupPath(os.Getenv("PWD"))
	// look for kool.yml on kool folder within user home directory
	_ = kool.AddLookupPath(path.Join(os.Getenv("HOME"), "kool"))

	script = args[0]

	if commands, err = kool.Parse(script); err != nil {
		if parser.IsMultipleDefinedScriptError(err) {
			// we should just warn the user about multiple finds for the script
			shell.Warning("Attention: the script was found in more than one kool.yml file")
		} else {
			shell.Error("failed parsing script")
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}

	if len(args) > 1 && len(commands) > 1 {
		shell.Error("error: you cannot pass in extra arguments to multiple commands scripts")
		os.Exit(2)
	}

	for _, command := range commands {
		if len(args) > 0 {
			command.AppendArgs(args...)
		}

		err = command.Interactive()

		if err != nil {
			shell.Error(err)
			os.Exit(1)
		}
	}
}
