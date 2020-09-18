package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// NewInfoCmd initializes new kool info command
func NewInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Prints out information about kool setup (like environment variables)",
		Run:   runInfo,
		Args:  cobra.MaximumNArgs(1),
	}
}

func init() {
	rootCmd.AddCommand(NewInfoCmd())
}

func runInfo(cmf *cobra.Command, args []string) {
	var filter string = "KOOL_"

	if len(args) > 0 {
		filter = args[0]
	}

	for _, envVar := range os.Environ() {
		if strings.Contains(envVar, filter) {
			fmt.Fprintf(cmf.OutOrStdout(), "%s\n", envVar)
		}
	}
}
