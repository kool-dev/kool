package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fwd",
	Short: "fwd - fwd - Fucking aWesome Development",
	Long: `An easy and robust software development environment
					tool helping you from project creation until deployment.
					Complete documentation is available at https://fwd.tools`,
}

// Execute proxies the call to cobra root command
func Execute() error {
	return rootCmd.Execute()
}
