package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// DockerFlags holds the flags for the docker command
type DockerFlags struct {
	DisableTty bool
}

var dockerCmd = &cobra.Command{
	Use:                "docker [image] [command]",
	Short:              "Creates a new container and runs the command in it.",
	Args:               cobra.MinimumNArgs(1),
	Run:                runDocker,
	DisableFlagParsing: true,
}

var dockerFlags = &DockerFlags{false}

func init() {
	rootCmd.AddCommand(dockerCmd)

	dockerCmd.Flags().BoolVarP(&dockerFlags.DisableTty, "disable-tty", "T", false, "Disables TTY")
}

func runDocker(cmd *cobra.Command, originalArgs []string) {
	var (
		args []string
	)

	if originalArgs[0] == "--disable-tty" || originalArgs[0] == "-T" {
		dockerFlags.DisableTty = true
		args = originalArgs[1:]
	} else {
		args = originalArgs
	}

	execDockerRun(args)
}

func execDockerRun(command []string) {
	var (
		args    []string
		err     error
		workDir string
		image   string
	)

	workDir, err = os.Getwd()
	args = []string{"run", "--init", "--rm", "-w", "/app"}
	if disableTty := os.Getenv("KOOL_TTY_DISABLE"); !dockerFlags.DisableTty && !(disableTty == "1" || disableTty == "true") {
		args = append(args, "-ti")
	}

	for i, arg := range command {
		if !strings.HasPrefix(arg, "-") {
			// we found the image!
			image = arg

			// this means the expected command is only the remaining items
			// so we override it
			command = command[i+1:]
			break
		}

		// didn't find the image (first non-option argument)
		// so we assume it is a pre-image option for docker runs
		args = append(args, arg)
	}

	if asuser := os.Getenv("KOOL_ASUSER"); asuser != "" && (strings.HasPrefix(image, "kooldev") || strings.HasPrefix(image, "fireworkweb")) {
		args = append(args, "--env", "ASUSER="+os.Getenv("KOOL_ASUSER"))
	}
	args = append(args, "--volume", workDir+":/app", image)
	args = append(args, command...)

	err = shellInteractive("docker", args...)

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
}
