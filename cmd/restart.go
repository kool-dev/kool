package cmd

import (
	"github.com/spf13/cobra"
)

// KoolRestart holds data and logic necessary for
// restarting containers controller by kool.
type KoolRestart struct {
	DefaultKoolService

	start KoolService
	stop  KoolService
}

// Execute implements the logic for restarting
func (r *KoolRestart) Execute(args []string) (err error) {
	if err = r.stop.Execute(nil); err != nil {
		return
	}

	err = r.start.Execute(nil)
	return
}

// NewRestartCommand initializes new kool start command
func NewRestartCommand(restart *KoolRestart) *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart containers - the same as stop followed by start.",
		Run: func(cmd *cobra.Command, args []string) {
			restart.SetWriter(cmd.OutOrStdout())
			restart.start.SetWriter(cmd.OutOrStdout())
			restart.stop.SetWriter(cmd.OutOrStdout())

			if err := restart.Execute(nil); err != nil {
				restart.Error(err)
				restart.Exit(1)
			}
		},
		DisableFlagsInUseLine: true,
	}
}

func init() {
	rootCmd.AddCommand(NewRestartCommand(&KoolRestart{
		*newDefaultKoolService(),
		NewKoolStart(),
		NewKoolStop(),
	}))
}
