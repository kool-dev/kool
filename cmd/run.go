package cmd

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/parser"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// KoolRun holds handlers and functions to implement the run command logic
type KoolRun struct {
	DefaultKoolService
	parser       parser.Parser
	env          environment.EnvStorage
	promptSelect shell.PromptSelect
	commands     []builder.Command
}

// ErrExtraArguments Extra arguments error
var ErrExtraArguments = errors.New("error: you cannot pass in extra arguments to multiple commands scripts")

// ErrKoolScriptNotFound means that the given script was not found
var ErrKoolScriptNotFound = errors.New("script was not found in any kool.yml file")

func init() {
	var (
		run    = NewKoolRun()
		runCmd = NewRunCommand(run)
	)

	rootCmd.AddCommand(runCmd)

	SetRunUsageFunc(run, runCmd)
}

// NewKoolRun creates a new handler for run logic with default dependencies
func NewKoolRun() *KoolRun {
	return &KoolRun{
		*newDefaultKoolService(),
		parser.NewParser(),
		environment.NewEnvStorage(),
		shell.NewPromptSelect(),
		[]builder.Command{},
	}
}

// Execute runs the run logic with incoming arguments.
func (r *KoolRun) Execute(originalArgs []string) (err error) {
	var (
		script string   = originalArgs[0]
		args   []string = originalArgs[1:]
	)

	// look for kool.yml on current working directory
	_ = r.parser.AddLookupPath(r.env.Get("PWD"))
	// look for kool.yml on kool folder within user home directory
	_ = r.parser.AddLookupPath(path.Join(r.env.Get("HOME"), "kool"))

	if r.commands, err = parseScript(r, script); err != nil {
		if parser.IsMultipleDefinedScriptError(err) {
			r.Warning("Attention: the script was found in more than one kool.yml file")
			err = nil
		} else {
			return
		}
	}

	if len(r.commands) == 0 {
		err = ErrKoolScriptNotFound
		return
	}

	if len(args) > 0 && len(r.commands) > 1 {
		err = ErrExtraArguments
		return
	}

	for _, command := range r.commands {
		if len(args) > 0 {
			command.AppendArgs(args...)
		}

		if err = r.Interactive(command); err != nil {
			return
		}
	}
	return
}

// NewRunCommand initializes new kool stop command
func NewRunCommand(run *KoolRun) (runCmd *cobra.Command) {
	runCmd = &cobra.Command{
		Use:   "run [SCRIPT]",
		Short: "Runs a custom command defined at kool.yaml in the working directory or in the kool folder of the user's home directory",
		Args:  cobra.MinimumNArgs(1),
		Run:   DefaultCommandRunFunction(run),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			return compListScripts(toComplete, run), cobra.ShellCompDirectiveNoFileComp
		},
	}

	// after a non-flag arg, stop parsing flags
	runCmd.Flags().SetInterspersed(false)

	return
}

// SetRunUsageFunc overrides usage function
func SetRunUsageFunc(run *KoolRun, runCmd *cobra.Command) {
	originalUsageText := runCmd.UsageString()
	runCmd.SetUsageFunc(getRunUsageFunc(run, originalUsageText))
}

func getRunUsageFunc(run *KoolRun, originalUsageText string) func(*cobra.Command) error {
	return func(cmd *cobra.Command) (err error) {
		var (
			sb      strings.Builder
			scripts []string
		)

		// look for kool.yml on current working directory
		_ = run.parser.AddLookupPath(run.env.Get("PWD"))
		// look for kool.yml on kool folder within user home directory
		_ = run.parser.AddLookupPath(path.Join(run.env.Get("HOME"), "kool"))

		if scripts, err = run.parser.ParseAvailableScripts(""); err != nil {
			if run.env.IsTrue("KOOL_VERBOSE") {
				run.Println("$ got an error trying to add available scripts to command usage template; error:", err.Error())
			}
			return
		}

		sb.WriteString(originalUsageText)
		sb.WriteString("\n")
		sb.WriteString("Available Scripts:\n")

		for _, script := range scripts {
			sb.WriteString("  ")
			sb.WriteString(script)
			sb.WriteString("\n")
		}

		run.Println(sb.String())
		return
	}
}

func compListScripts(toComplete string, run *KoolRun) (scripts []string) {
	var err error
	// look for kool.yml on current working directory
	_ = run.parser.AddLookupPath(run.env.Get("PWD"))
	// look for kool.yml on kool folder within user home directory
	_ = run.parser.AddLookupPath(path.Join(run.env.Get("HOME"), "kool"))

	if scripts, err = run.parser.ParseAvailableScripts(toComplete); err != nil {
		return nil
	}

	return
}

func parseScript(run *KoolRun, script string) (commands []builder.Command, err error) {
	var (
		similarIsCorrect string
		chosenSimilar    string
	)

	if commands, err = run.parser.Parse(script); err != nil {
		if parser.IsPossibleTypoError(err) && run.IsTerminal() {
			if similarIsCorrect, _ = run.promptSelect.Ask(err.Error(), []string{"Yes", "No"}); similarIsCorrect != "Yes" {
				err = ErrKoolScriptNotFound
				return
			}

			if possibleScripts := err.(*parser.ErrPossibleTypo).Similars(); len(possibleScripts) == 1 {
				chosenSimilar = possibleScripts[0]
			} else {
				chosenSimilar, _ = run.promptSelect.Ask("which one did you mean?", possibleScripts)
			}

			commands, err = run.parser.Parse(chosenSimilar)
		}
	}

	return
}
