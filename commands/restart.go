package commands

import (
	"github.com/spf13/cobra"
)

// KoolRestartFlags holds the flags for the kool restart command
type KoolRestartFlags struct {
	Purge   bool
	Rebuild bool
}

// NewRestartCommand initializes new kool start command
func NewRestartCommand(stop KoolService, start KoolService) (restartCmd *cobra.Command) {
	var flags *KoolRestartFlags = &KoolRestartFlags{false, false}

	restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart running service containers (the same as 'kool stop' followed by 'kool start')",
		Run: func(cmd *cobra.Command, args []string) {
			if _, ok := stop.(*KoolStop); ok && flags.Purge {
				stop.(*KoolStop).Flags.Purge = true
			}
			if _, ok := start.(*KoolStart); ok && flags.Rebuild {
				start.(*KoolStart).Flags.Rebuild = true
			}
			DefaultCommandRunFunction(stop, start)(cmd, args)
		},

		DisableFlagsInUseLine: true,
	}

	restartCmd.Flags().BoolVarP(&flags.Purge, "purge", "", false, "Remove all persistent data from volume mounts on containers")
	restartCmd.Flags().BoolVarP(&flags.Rebuild, "rebuild", "", false, "Updates and builds service's images")

	return
}

func AddKoolRestart(root *cobra.Command) {
	root.AddCommand(NewRestartCommand(NewKoolStop(), NewKoolStart()))
}
