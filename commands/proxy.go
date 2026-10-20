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
	proxyCommand.AddCommand(newProxyAliasesCommand(environment.NewEnvStorage()))
	root.AddCommand(proxyCommand)
}

func newProxyAliasesCommand(env environment.EnvStorage) *cobra.Command {
	return &cobra.Command{
		Use:   "aliases",
		Short: "List configured internal proxy aliases for the current project",
		Long:  "List configured internal Docker DNS aliases for the current project. This does not inspect running containers or contact Caddy.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			aliases, err := proxy.NewManager(shell.NewShell(), env).Aliases()
			if err != nil {
				return err
			}
			if len(aliases) == 0 {
				cmd.Println("No proxy aliases configured.")
				return nil
			}
			table := shell.NewTableWriter()
			table.SetWriter(cmd.OutOrStdout())
			table.AppendHeader("Project", "Service", "Network", "Alias")
			for _, alias := range aliases {
				table.AppendRow(alias.Project, alias.Service, alias.Network, alias.Name)
			}
			table.Render()
			return nil
		},
	}
}
