package commands

import (
	"github.com/spf13/cobra"
)

// KoolRestartFlags holds the flags for the kool restart command
type KoolRestartFlags struct {
	Purge bool
}

var flags *KoolRestartFlags = &KoolRestartFlags{false}

// NewRestartCommand initializes new kool start command
func NewRestartCommand(stop KoolService, start KoolService) (restartCmd *cobra.Command) {
	restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart running service containers (the same as 'kool stop' followed by 'kool start')",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := stop.(*KoolStop); ok && flags.Purge {
				stop.(*KoolStop).Flags.Purge = true
			}
			DefaultCommandRunFunction(stop, start)(cmd, args)
		},

		DisableFlagsInUseLine: true,
	}

	restartCmd.Flags().BoolVarP(&flags.Purge, "purge", "", false, "Remove all persistent data from volume mounts on containers")

	return
}

func AddKoolRestart(root *cobra.Command) {
	root.AddCommand(NewRestartCommand(NewKoolStop(), NewKoolStart()))
}
