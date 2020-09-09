package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// DockerFlags holds the flags for the docker command
type DockerFlags struct {
	DisableTty bool
	EnvVariables []string
}

var dockerCmd = &cobra.Command{
	Use:                "docker [options] [image] [command]",
	Args:               cobra.MinimumNArgs(1),
	Run:                runDocker,
	Short:              "Creates a new container and runs the command in it.",
	Long: `This command acts as a helper for docker run.
You can start with options that go before the image name
for docker run itself, i.e --env='VAR=VALUE'. Then you must pass
the image name and the command you want to exucute on that image.`,
}

var dockerFlags = &DockerFlags{false, []string{}}

func init() {
	rootCmd.AddCommand(dockerCmd)

	dockerCmd.Flags().BoolVarP(&dockerFlags.DisableTty, "disable-tty", "T", false, "Disables TTY")
	dockerCmd.Flags().StringArrayVarP(&dockerFlags.EnvVariables, "env", "e", []string{}, "Environment variables")

	//After a non-flag arg, stop parsing flags
	dockerCmd.Flags().SetInterspersed(false)
}

func runDocker(docker *cobra.Command, args []string) {
	image := args[0]
	command := args[1:]

	execDockerRun(image, command)
}

func execDockerRun(image string, command []string) {
	var (
		args    []string
		err     error
		workDir string
	)

	workDir, err = os.Getwd()
	args = []string{"run", "--init", "--rm", "-w", "/app"}
	if disableTty := os.Getenv("KOOL_TTY_DISABLE"); !dockerFlags.DisableTty && !(disableTty == "1" || disableTty == "true") {
		args = append(args, "-ti")
	}

	if asuser := os.Getenv("KOOL_ASUSER"); asuser != "" && (strings.HasPrefix(image, "kooldev") || strings.HasPrefix(image, "fireworkweb")) {
		args = append(args, "--env", "ASUSER="+os.Getenv("KOOL_ASUSER"))
	}

	if len(dockerFlags.EnvVariables) > 0 {
		for _, envVar := range dockerFlags.EnvVariables {
			args = append(args, "--env", envVar)
		}
	}

	args = append(args, "--volume", workDir+":/app", image)
	args = append(args, command...)

	err = shellInteractive("docker", args...)

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
}
