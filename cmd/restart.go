package cmd

import (
	"github.com/spf13/cobra"
)

// NewRestartCommand initializes new kool start command
func NewRestartCommand(stop KoolService, start KoolService) *cobra.Command {
	return &cobra.Command{
		Use:                   "restart",
		Short:                 "Restart containers - the same as stop followed by start.",
		Run:                   DefaultCommandRunFunction(stop, start),
		DisableFlagsInUseLine: true,
	}
}

func init() {
	rootCmd.AddCommand(NewRestartCommand(NewKoolStop(), NewKoolStart()))
}
