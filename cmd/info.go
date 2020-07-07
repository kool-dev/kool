package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints out information about kool setup (like environment variables)",
	Run:   runInfo,
	Args:  cobra.MaximumNArgs(1),
}

func runInfo(cmf *cobra.Command, args []string) {
	var filter string = "KOOL_"

	if len(args) > 0 {
		filter = args[0]
	}

	for _, envVar := range os.Environ() {
		if strings.Contains(envVar, filter) {
			fmt.Println(envVar)
		}
	}
}
