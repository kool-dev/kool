package commands

import (
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// KoolDockerFlags holds the flags for the docker command
type KoolDockerFlags struct {
	EnvVariables []string
	Volumes      []string
	Publish      []string
	Network      []string
}

// KoolDocker holds handlers and functions to implement the docker command logic
type KoolDocker struct {
	DefaultKoolService
	Flags *KoolDockerFlags

	envStorage environment.EnvStorage
	dockerRun  builder.Command
}

func AddKoolDocker(root *cobra.Command) {
	var (
		docker    = NewKoolDocker()
		dockerCmd = NewDockerCommand(docker)
	)

	root.AddCommand(dockerCmd)
}

// NewKoolDocker creates a new handler for docker logic
func NewKoolDocker() *KoolDocker {
	return &KoolDocker{
		*newDefaultKoolService(),
		&KoolDockerFlags{[]string{}, []string{}, []string{}, []string{}},
		environment.NewEnvStorage(),
		builder.NewCommand("docker", "run", "--init", "--rm", "-w", "/app", "-i"),
	}
}

// Execute runs the docker logic with incoming arguments.
func (d *KoolDocker) Execute(args []string) (err error) {
	workDir, _ := os.Getwd()

	if d.Shell().IsTerminal() {
		d.dockerRun.AppendArgs("-t")
	}

	// only adds env ASUSER if we are not running on MacOS
	if asuser := d.envStorage.Get("KOOL_ASUSER"); asuser != "" && runtime.GOOS != "darwin" {
		d.dockerRun.AppendArgs("--env", "ASUSER="+asuser)
	}

	if len(d.Flags.EnvVariables) > 0 {
		for _, envVar := range d.Flags.EnvVariables {
			d.dockerRun.AppendArgs("--env", envVar)
		}
	}

	d.dockerRun.AppendArgs("--volume", workDir+":/app:delegated")

	if len(d.Flags.Volumes) > 0 {
		for _, volume := range d.Flags.Volumes {
			d.dockerRun.AppendArgs("--volume", volume)
		}
	}

	if len(d.Flags.Publish) > 0 {
		for _, publish := range d.Flags.Publish {
			d.dockerRun.AppendArgs("--publish", publish)
		}
	}

	if len(d.Flags.Network) > 0 {
		for _, network := range d.Flags.Network {
			d.dockerRun.AppendArgs("--network", network)
		}
	}

	err = d.Shell().Interactive(d.dockerRun, args...)
	return
}

// NewDockerCommand initializes new kool docker command
func NewDockerCommand(docker *KoolDocker) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "docker [OPTIONS] IMAGE [COMMAND] [--] [ARG...]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Create a new container (a powered up 'docker run')",
		Long: `A helper for 'docker run'. Any [OPTIONS] added before the
IMAGE name will be used by 'docker run' itself (i.e. --env='VAR=VALUE').
Add an optional [COMMAND] to execute on the IMAGE, and use [--] after
the [COMMAND] to provide optional arguments required by the COMMAND.`,
		RunE: DefaultCommandRunFunction(docker),

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().StringArrayVarP(&docker.Flags.EnvVariables, "env", "e", []string{}, "Environment variables.")
	cmd.Flags().StringArrayVarP(&docker.Flags.Volumes, "volume", "v", []string{}, "Bind mount a volume.")
	cmd.Flags().StringArrayVarP(&docker.Flags.Publish, "publish", "p", []string{}, "Publish a container's port(s) to the host.")
	cmd.Flags().StringArrayVarP(&docker.Flags.Network, "network", "n", []string{}, "Connect a container to a network.")

	//After a non-flag arg, stop parsing flags
	cmd.Flags().SetInterspersed(false)

	return
}
