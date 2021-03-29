package cmd

import "github.com/spf13/cobra"

// KoolCompletion holds handlers and functions to implement the docker command logic
type KoolCompletion struct {
	DefaultKoolService
	rootCmd *cobra.Command
}

func AddKoolCompletion(root *cobra.Command) {
	var (
		completion    = NewKoolCompletion(root)
		completionCmd = NewCompletionCommand(completion)
	)

	root.AddCommand(completionCmd)
}

// NewKoolCompletion creates a new handler for completion logic
func NewKoolCompletion(root *cobra.Command) *KoolCompletion {
	return &KoolCompletion{
		*newDefaultKoolService(),
		root,
	}
}

// Execute runs the completion logic with incoming arguments.
func (c *KoolCompletion) Execute(args []string) (err error) {
	switch args[0] {
	case "bash":
		err = c.rootCmd.GenBashCompletion(c.OutStream())
	case "zsh":
		err = c.rootCmd.GenZshCompletion(c.OutStream())
	case "fish":
		err = c.rootCmd.GenFishCompletion(c.OutStream(), true)
	case "powershell":
		err = c.rootCmd.GenPowerShellCompletion(c.OutStream())
	}
	return
}

// NewCompletionCommand initializes new kool completion command
func NewCompletionCommand(completion *KoolCompletion) *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion configuration script.",
		Long: `Autocompletion:

If you want to use kool autocompletion in your Unix shell, follow the appropriate instructions below.

After running one of the below commands, remember to start a new shell for autocompletion to take effect.

#### Bash

Temporarily enable autocompletion for your current session only:

$ source <(kool completion bash)

Permanently enable autocompletion for all sessions:

  Linux:

  $ kool completion bash > /etc/bash_completion.d/kool

  macOS:

  $ kool completion bash > /usr/local/etc/bash_completion.d/kool

#### Zsh

If Zsh tab completion is not already initialized on your machine, run the following command to turn it on.

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

Permanently enable autocompletion for all sessions:

$ kool completion zsh > "${fpath[1]}/_kool"

#### Fish

Temporarily enable autocompletion for your current session only:

$ kool completion fish | source

Permanently enable autocompletion for all sessions:

$ kool completion fish > ~/.config/fish/completions/kool.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Hidden:                true,
		Run:                   DefaultCommandRunFunction(completion),
	}
}
