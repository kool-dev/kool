package cmd

import "github.com/spf13/cobra"

// KoolCompletion holds handlers and functions to implement the docker command logic
type KoolCompletion struct {
	DefaultKoolService
	rootCmd *cobra.Command
}

func init() {
	var (
		completion    = NewKoolCompletion()
		completionCmd = NewCompletionCommand(completion)
	)

	rootCmd.AddCommand(completionCmd)
}

// NewKoolCompletion creates a new handler for completion logic
func NewKoolCompletion() *KoolCompletion {
	return &KoolCompletion{
		*newDefaultKoolService(),
		rootCmd,
	}
}

// Execute runs the completion logic with incoming arguments.
func (c *KoolCompletion) Execute(args []string) (err error) {
	switch args[0] {
	case "bash":
		err = c.rootCmd.GenBashCompletion(c.GetWriter())
	case "zsh":
		err = c.rootCmd.GenZshCompletion(c.GetWriter())
	case "fish":
		err = c.rootCmd.GenFishCompletion(c.GetWriter(), true)
	case "powershell":
		err = c.rootCmd.GenPowerShellCompletion(c.GetWriter())
	}
	return
}

// NewCompletionCommand initializes new kool completion command
func NewCompletionCommand(completion *KoolCompletion) *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(kool completion bash)

# To load completions for each session, execute once:
Linux:
  $ kool completion bash > /etc/bash_completion.d/kool
MacOS:
  $ kool completion bash > /usr/local/etc/bash_completion.d/kool

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ kool completion zsh > "${fpath[1]}/_kool"

# You will need to start a new shell for this setup to take effect.

Fish:

$ kool completion fish | source

# To load completions for each session, execute once:
$ kool completion fish > ~/.config/fish/completions/kool.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Hidden:                true,
		Run: func(cmd *cobra.Command, args []string) {
			completion.SetWriter(cmd.OutOrStdout())

			if err := completion.Execute(args); err != nil {
				completion.Error(err)
				completion.Exit(1)
			}
		},
	}
}
