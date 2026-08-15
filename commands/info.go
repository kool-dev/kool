package commands

import (
	"encoding/json"
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

	envStorage                  environment.EnvStorage
	cmdDocker, cmdDockerCompose builder.Command
}

type infoOutputJSON struct {
	KoolVersion          string            `json:"kool_version"`
	KoolBinPath          string            `json:"kool_bin_path"`
	DockerVersion        string            `json:"docker_version"`
	DockerBinPath        string            `json:"docker_bin_path"`
	DockerComposeVersion string            `json:"docker_compose_version"`
	Env                  map[string]string `json:"env"`
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
		builder.NewCommand("docker", "compose", "version"),
	}
}

func AddKoolInfo(root *cobra.Command) {
	root.AddCommand(NewInfoCmd(NewKoolInfo()))
}

// Execute executes info logic
func (i *KoolInfo) Execute(args []string) (err error) {
	var (
		filter = "KOOL_"
		output string
	)

	if len(args) > 0 {
		filter = args[0]
	}

	if i.Shell().IsJSONOutput() {
		return i.executeJSON(filter)
	}

	// kool CLI info
	i.Shell().Println("Kool Version ", version)
	if output, err = os.Executable(); err != nil {
		return
	}
	i.Shell().Println("Kool Bin Path:", output)

	i.Shell().Println("")
	// docker CLI info
	if output, err = i.Shell().Exec(i.cmdDocker); err != nil {
		return
	}
	i.Shell().Println(output)

	if err = i.shell.LookPath(i.cmdDocker); err != nil {
		return
	}
	output, _ = exec.LookPath(i.cmdDocker.Cmd())

	i.Shell().Println("Docker Bin Path:", output)

	i.Shell().Println("")

	// docker compose v2 info
	if output, err = i.Shell().Exec(i.cmdDockerCompose); err != nil {
		// just alert missing docker compose, but don't elevate error
		i.Shell().Warning("Docker Compose:", err.Error())
		i.Shell().Error(fmt.Errorf("you need to have Docker Compose V2 available; make sure to update your Docker installation"))
		return
	} else {
		i.Shell().Println(output)
	}

	i.Shell().Println("")
	i.Shell().Println("Environment Variables of Interest:")
	i.Shell().Println("")

	for _, envVar := range i.envStorage.All() {
		if strings.Contains(envVar, filter) {
			// keep from printing out known to be sensitive values
			if strings.Contains(envVar, "KOOL_API_TOKEN") {
				i.Shell().Warning("KOOL_API_TOKEN=***************** [redacted]")
			} else {
				i.Shell().Println(envVar)
			}
		}
	}

	i.Shell().Println("")
	i.Shell().Println("kool installation seems to be working as expected.")

	return
}

func (i *KoolInfo) executeJSON(filter string) (err error) {
	var (
		output string
		info   infoOutputJSON
	)

	info.KoolVersion = version

	if output, err = os.Executable(); err != nil {
		return
	}
	info.KoolBinPath = output

	if output, err = i.Shell().Exec(i.cmdDocker); err != nil {
		return
	}
	info.DockerVersion = output

	if err = i.shell.LookPath(i.cmdDocker); err != nil {
		return
	}
	info.DockerBinPath, _ = exec.LookPath(i.cmdDocker.Cmd())

	if output, err = i.Shell().Exec(i.cmdDockerCompose); err != nil {
		i.Shell().Warning("Docker Compose:", err.Error())
		i.Shell().Error(fmt.Errorf("you need to have Docker Compose V2 available; make sure to update your Docker installation"))
		return
	}
	info.DockerComposeVersion = output

	info.Env = map[string]string{}
	for _, envVar := range i.envStorage.All() {
		if !strings.Contains(envVar, filter) {
			continue
		}
		parts := strings.SplitN(envVar, "=", 2)
		key, value := parts[0], ""
		if len(parts) > 1 {
			value = parts[1]
		}
		if key == "KOOL_API_TOKEN" {
			value = "***************** [redacted]"
		}
		info.Env[key] = value
	}

	var payload []byte
	if payload, err = json.Marshal(info); err != nil {
		return
	}
	i.Shell().Println(string(payload))
	return
}
