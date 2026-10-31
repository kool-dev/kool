package commands

import (
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/proxy"

	"github.com/spf13/cobra"
)

// AddKoolProxy adds commands for managing Kool's local proxy.
func AddKoolProxy(root *cobra.Command) {
	proxyCommand := &cobra.Command{
		Use:   "proxy",
		Short: "Manage Kool's local proxy",
	}
	proxyCommand.AddCommand(&cobra.Command{
		Use:   "trust",
		Short: "Trust the local proxy certificate authority",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := shell.NewShell()
			sh.SetInStream(cmd.InOrStdin())
			sh.SetOutStream(cmd.OutOrStdout())
			sh.SetErrStream(cmd.ErrOrStderr())
			return proxy.NewManager(sh, environment.NewEnvStorage()).Trust()
		},
	})
	root.AddCommand(proxyCommand)
}
