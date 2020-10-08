package cmd

import (
	"kool-dev/kool/environment"
	"strings"

	"github.com/spf13/cobra"
)

// KoolInfo holds handlers and functions for info logic
type KoolInfo struct {
	DefaultKoolService

	envStorage environment.EnvStorage
}

// NewInfoCmd initializes new kool info command
func NewInfoCmd(info *KoolInfo) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Prints out information about kool setup (like environment variables)",
		Run:   DefaultCommandRunFunction(info),
		Args:  cobra.MaximumNArgs(1),
	}
}

// NewKoolInfo creates a new pointer with default KoolInfo service
func NewKoolInfo() *KoolInfo {
	return &KoolInfo{
		*newDefaultKoolService(),
		environment.NewEnvStorage(),
	}
}

func init() {
	rootCmd.AddCommand(NewInfoCmd(NewKoolInfo()))
}

// Execute executes info logic
func (i *KoolInfo) Execute(args []string) (err error) {
	var filter string = "KOOL_"

	if len(args) > 0 {
		filter = args[0]
	}

	for _, envVar := range i.envStorage.All() {
		if strings.Contains(envVar, filter) {
			i.Println(envVar)
		}
	}
	return
}
