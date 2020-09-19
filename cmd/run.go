package cmd

import (
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

var runsCmdOutputWriter shell.OutputWriter = shell.NewOutputWriter()

func init() {
	rootCmd.AddCommand(runCmd)

	// after a non-flag arg, stop parsing flags
	runCmd.Flags().SetInterspersed(false)
}

func runRun(cmd *cobra.Command, originalArgs []string) {
	var (
		script   string
		args     []string
		commands []*builder.DefaultCommand
		err      error
		kool     parser.Parser
	)

	runsCmdOutputWriter.SetWriter(cmd.OutOrStdout())

	kool = parser.NewParser()

	// look for kool.yml on current working directory
	_ = kool.AddLookupPath(os.Getenv("PWD"))
	// look for kool.yml on kool folder within user home directory
	_ = kool.AddLookupPath(path.Join(os.Getenv("HOME"), "kool"))

	script = originalArgs[0]
	args = originalArgs[1:]

	if commands, err = kool.Parse(script); err != nil {
		if parser.IsMultipleDefinedScriptError(err) {
			// we should just warn the user about multiple finds for the script
			runsCmdOutputWriter.Warning("Attention: the script was found in more than one kool.yml file")
		} else {
			runsCmdOutputWriter.ExecError("failed parsing script", err)
			os.Exit(1)
		}
	}

	if len(args) > 0 && len(commands) > 1 {
		runsCmdOutputWriter.Error("error: you cannot pass in extra arguments to multiple commands scripts")
		os.Exit(2)
	}

	for _, command := range commands {
		if len(args) > 0 {
			command.AppendArgs(args...)
		}

		err = command.Interactive()

		if err != nil {
			runsCmdOutputWriter.Error(err)
			os.Exit(1)
		}
	}
}
