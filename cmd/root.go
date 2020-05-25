package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fwd",
	Short: "fwd - FWD Web Development - easy environments",
	Long: `An easy and robust software development environment
					tool helping you from project creation until deployment.
					Complete documentation is available at https://fwd.tools`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fwd - Fucking aWesome Development")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
