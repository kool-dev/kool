package commands

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// KoolInfo holds handlers and functions for info logic
type KoolInfo struct {
	DefaultKoolService

	envStorage environment.EnvStorage
	cmdDocker, cmdDockerCompose builder.Command
}

// NewInfoCmd initializes new kool info command
func NewInfoCmd(info *KoolInfo) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Print out information about the local environment",
		Long:  "Print out information about the local environment, such as environment variables.",
		RunE:  DefaultCommandRunFunction(info),
		Args:  cobra.MaximumNArgs(1),

		DisableFlagsInUseLine: true,
	}
}

// NewKoolInfo creates a new pointer with default KoolInfo service
func NewKoolInfo() *KoolInfo {
	return &KoolInfo{
		*newDefaultKoolService(),
		environment.NewEnvStorage(),
		builder.NewCommand("docker", "-v"),
		builder.NewCommand("docker-compose", "-v"),
	}
}

func AddKoolInfo(root *cobra.Command) {
	root.AddCommand(NewInfoCmd(NewKoolInfo()))
}

// Execute executes info logic
func (i *KoolInfo) Execute(args []string) (err error) {
	var (
		filter string = "KOOL_"
		output string
	)

	if len(args) > 0 {
		filter = args[0]
	}

	// kool CLI info
	i.Println("Kool Version ", version)
	if output, err = os.Executable(); err != nil {
		fmt.Println("err1")
		return
	}
	i.Println("Kool Bin Path:", output)

	i.Println("")
	// docker CLI info
	if output, err = i.Exec(i.cmdDocker); err != nil {
		return
	}
	i.Println(output)

	if err = i.shell.LookPath(i.cmdDocker); err != nil {
		return
	}
	output, _ = exec.LookPath(i.cmdDocker.Cmd())

	i.Println("Docker Bin Path:", output)

	i.Println("")

	// docker-compose CLI info
	if output, err = i.Exec(i.cmdDockerCompose); err != nil {
		// just alert missing docker-compose, but don't elevate error
		i.Warning("Docker Compose:", err.Error())
		i.Warning("It's okay not having docker-compose installed, as kool will fallback to using a container for it when necessary.")
		err = nil
	} else {
		i.Println(output)
		output, _ = exec.LookPath("docker-compose")
		i.Println("Docker Compose Bin Path:", output)
	}

	i.Println("")
	i.Println("Environment Variables of Interest:")
	i.Println("")

	for _, envVar := range i.envStorage.All() {
		if strings.Contains(envVar, filter) {
			// keep from printing out known to be sensitive values
			if strings.Contains(envVar, "KOOL_API_TOKEN") {
				i.Warning("KOOL_API_TOKEN=***************** [redacted]")
			} else {
				i.Println(envVar)
			}
		}
	}

	return
}
