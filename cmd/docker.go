package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"os"

	"github.com/spf13/cobra"
)

// KoolDockerFlags holds the flags for the docker command
type KoolDockerFlags struct {
	DisableTty   bool
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
		&KoolDockerFlags{false, []string{}, []string{}, []string{}, []string{}},
		environment.NewEnvStorage(),
		builder.NewCommand("docker", "run", "--init", "--rm", "-w", "/app", "-i"),
	}
}

// Execute runs the docker logic with incoming arguments.
func (d *KoolDocker) Execute(args []string) (err error) {
	workDir, _ := os.Getwd()

	if d.IsTerminal() {
		d.dockerRun.AppendArgs("-t")
	}

	if asuser := d.envStorage.Get("KOOL_ASUSER"); asuser != "" {
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

	err = d.Interactive(d.dockerRun, args...)
	return
}

// NewDockerCommand initializes new kool docker command
func NewDockerCommand(docker *KoolDocker) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "docker [OPTIONS] IMAGE [COMMAND] -- [ARG...]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Create a new container (a powered up 'docker run')",
		Long: `This command acts as a helper for 'docker run'.
You can provide one or more [OPTIONS] before the IMAGE name that will be used
by 'docker run' itself (i.e --env='VAR=VALUE'). Then you must pass
the IMAGE name and the [COMMAND] you want to execute on that IMAGE. After that you can use -- and follow with any extra arguments that command may require.`,
		Run: DefaultCommandRunFunction(docker),

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().BoolVarP(&docker.Flags.DisableTty, "disable-tty", "T", false, "Deprecated - no effect.")
	cmd.Flags().StringArrayVarP(&docker.Flags.EnvVariables, "env", "e", []string{}, "Environment variables.")
	cmd.Flags().StringArrayVarP(&docker.Flags.Volumes, "volume", "v", []string{}, "Bind mount a volume.")
	cmd.Flags().StringArrayVarP(&docker.Flags.Publish, "publish", "p", []string{}, "Publish a container's port(s) to the host.")
	cmd.Flags().StringArrayVarP(&docker.Flags.Network, "network", "n", []string{}, "Connect a container to a network.")

	//After a non-flag arg, stop parsing flags
	cmd.Flags().SetInterspersed(false)

	return
}
