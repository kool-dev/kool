package cmd

import (
	"github.com/spf13/cobra"
)

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
