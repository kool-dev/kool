package cmd

import (
	"github.com/spf13/cobra"
)

type cobraRunFN func(*cobra.Command, []string)

// NewCommandData holds data to create a new command
type NewCommandData struct {
	Use, Short, Long string
	Run              cobraRunFN
}

var version string = "0.0.0-dev"

var rootCmd = &cobra.Command{
	Use:   "kool",
	Short: "kool - Kool stuff",
	Long: `An easy and robust software development environment
tool helping you from project creation until deployment.
Complete documentation is available at https://kool.dev/docs`,
	Version:           version,
	DisableAutoGenTag: true,
}

// Execute proxies the call to cobra root command
func Execute() error {
	return rootCmd.Execute()
}

// RootCmd exposes the root command
func RootCmd() *cobra.Command {
	return rootCmd
}

// CreateCommand creates a new command
func CreateCommand(service KoolService, values NewCommandData) *cobra.Command {
	newCmd := &cobra.Command{
		Use:   values.Use,
		Short: values.Short,
		Long:  values.Long,
		PreRun: func(cmd *cobra.Command, args []string) {
			service.SetWriter(cmd.OutOrStdout())
		},
		Run: values.Run,
	}

	if newCmd.Run == nil {
		newCmd.Run = func(cmd *cobra.Command, args []string) {
			if err := service.Execute(args); err != nil {
				service.Error(err)
				service.Exit(1)
			}
		}
	}

	return newCmd
}
