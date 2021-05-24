package commands

import (
	"errors"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/parser"
	"kool-dev/kool/core/shell"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// KoolRunFlags holds the flags for the run command
type KoolRunFlags struct {
	EnvVariables []string
}

// KoolRun holds handlers and functions to implement the run command logic
type KoolRun struct {
	DefaultKoolService
	Flags        *KoolRunFlags
	parser       parser.Parser
	env          environment.EnvStorage
	promptSelect shell.PromptSelect
	commands     []builder.Command
}

// ErrExtraArguments Extra arguments error
var ErrExtraArguments = errors.New("error: you cannot pass in extra arguments to multiple commands scripts")

// ErrKoolScriptNotFound means that the given script was not found
var ErrKoolScriptNotFound = errors.New("script was not found in any kool.yml file")

func AddKoolRun(root *cobra.Command) {
	var (
		run    = NewKoolRun()
		runCmd = NewRunCommand(run)
	)

	root.AddCommand(runCmd)

	SetRunUsageFunc(run, runCmd)
}

// NewKoolRun creates a new handler for run logic with default dependencies
func NewKoolRun() *KoolRun {
	return &KoolRun{
		*newDefaultKoolService(),
		&KoolRunFlags{[]string{}},
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

	if err = r.parseScript(script); err != nil {
		return
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
		Use:   "run SCRIPT [--] [ARG...]",
		Short: "Execute a script defined in kool.yml",
		Long: `Execute the specified SCRIPT, as defined in the kool.yml file.
A single-line SCRIPT can be run with optional arguments.`,
		Args: cobra.MinimumNArgs(1),
		Run:  DefaultCommandRunFunction(run),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			return compListScripts(toComplete, run), cobra.ShellCompDirectiveNoFileComp
		},
		DisableFlagsInUseLine: true,
	}

	runCmd.Flags().StringArrayVarP(&run.Flags.EnvVariables, "env", "e", []string{}, "Environment variables.")

	// after a non-flag arg, stop parsing flags
	runCmd.Flags().SetInterspersed(false)

	return
}

// SetRunUsageFunc overrides usage function
func SetRunUsageFunc(run *KoolRun, runCmd *cobra.Command) {
	originalUsageText := runCmd.UsageString()
	runCmd.SetUsageFunc(getRunUsageFunc(run, originalUsageText))
}

func (r *KoolRun) parseScript(script string) (err error) {
	var (
		originalEnvs     map[string]string = make(map[string]string)
		similarIsCorrect string
		chosenSimilar    string
	)

	for _, envVar := range r.Flags.EnvVariables {
		pair := strings.SplitN(envVar, "=", 2)
		originalEnvs[pair[0]] = r.env.Get(pair[0])
		r.env.Set(pair[0], pair[1])
	}

	defer func() {
		for k, v := range originalEnvs {
			r.env.Set(k, v)
		}
	}()

	if r.commands, err = r.parser.Parse(script); err != nil {
		if parser.IsPossibleTypoError(err) && r.IsTerminal() {
			var promptError error

			similarIsCorrect, promptError = r.promptSelect.Ask(err.Error(), []string{"Yes", "No"})

			if promptError != nil {
				err = promptError
				return
			}

			if similarIsCorrect != "Yes" {
				err = ErrKoolScriptNotFound
				return
			}

			if possibleScripts := err.(*parser.ErrPossibleTypo).Similars(); len(possibleScripts) == 1 {
				chosenSimilar = possibleScripts[0]
			} else {
				chosenSimilar, promptError = r.promptSelect.Ask("which one did you mean?", possibleScripts)

				if promptError != nil {
					err = promptError
					return
				}
			}

			r.commands, err = r.parser.Parse(chosenSimilar)
			return
		}

		if parser.IsMultipleDefinedScriptError(err) {
			// we should just warn the user about multiple finds for the script
			r.Warning("Attention: the script was found in more than one kool.yml file")
			err = nil
		}
	}

	return
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
