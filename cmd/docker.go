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
	Volumes []string
	Publish []string
}

var dockerCmd = &cobra.Command{
	Use:                "docker [options] [image] [command]",
	Args:               cobra.MinimumNArgs(1),
	Run:                runDocker,
	Short:              "Creates a new container and runs the command in it.",
	Long: `This command acts as a helper for docker run.
You can start with options that go before the image name
for docker run itself, i.e --env='VAR=VALUE'. Then you must pass
the image name and the command you want to execute on that image.`,
}

var dockerFlags = &DockerFlags{false, []string{}, []string{}, []string{}}

func init() {
	rootCmd.AddCommand(dockerCmd)

	dockerCmd.Flags().BoolVarP(&dockerFlags.DisableTty, "disable-tty", "T", false, "Disables TTY")
	dockerCmd.Flags().StringArrayVarP(&dockerFlags.EnvVariables, "env", "e", []string{}, "Environment variables")
	dockerCmd.Flags().StringArrayVarP(&dockerFlags.Volumes, "volume", "v", []string{}, "Bind mount a volume")
	dockerCmd.Flags().StringArrayVarP(&dockerFlags.Publish, "publish", "p", []string{}, "Publish a container’s port(s) to the host")

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

	args = append(args, "--volume", workDir+":/app:delegated")

	if len(dockerFlags.Volumes) > 0 {
		for _, volume := range dockerFlags.Volumes {
			args = append(args, "--volume", volume)
		}
	}

	if len(dockerFlags.Publish) > 0 {
		for _, publish := range dockerFlags.Publish {
			args = append(args, "--publish", publish)
		}
	}

	args = append(args, image)
	args = append(args, command...)

	err = shellInteractive("docker", args...)

	if err != nil {
		execError("", err)
		os.Exit(1)
	}
}
