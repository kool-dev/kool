package cmd

import (
	"kool-dev/kool/cmd/shell"

	"github.com/spf13/cobra"
)

// KoolRestart holds data and logic necessary for
// restarting containers controller by kool.
type KoolRestart struct {
	start *KoolStart
	stop  *KoolStop

	exiter shell.Exiter
	out    shell.OutputWriter
}

// NewRestartCommand initializes new kool start command
func NewRestartCommand(restart *KoolRestart) *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart containers - the same as stop followed by start.",
		Run: func(cmd *cobra.Command, args []string) {
			restart.start.out.SetWriter(cmd.OutOrStdout())
			restart.stop.out.SetWriter(cmd.OutOrStdout())

			if err := restart.stop.Execute(nil); err != nil {
				restart.out.Error(err)
				restart.exiter.Exit(1)
			}

			if err := restart.start.Execute(nil); err != nil {
				restart.out.Error(err)
				restart.exiter.Exit(1)
			}

		},
		DisableFlagsInUseLine: true,
	}
}

var restartCmd = NewRestartCommand(&KoolRestart{
	NewKoolStart(),
	NewKoolStop(),
	shell.NewExiter(),
	shell.NewOutputWriter(),
})

func init() {
	rootCmd.AddCommand(restartCmd)
}
