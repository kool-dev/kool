package cmd

import (
	"github.com/spf13/cobra"
)

// NewRestartCommand initializes new kool start command
func NewRestartCommand(stop KoolService, start KoolService) *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart the service containers (the same as 'stop' followed by 'start')",
		Run:   DefaultCommandRunFunction(stop, start),

		DisableFlagsInUseLine: true,
	}
}

func AddKoolRestart(root *cobra.Command) {
	root.AddCommand(NewRestartCommand(NewKoolStop(), NewKoolStart()))
}
