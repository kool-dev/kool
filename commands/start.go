package commands

import (
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/network"
	"kool-dev/kool/services/checker"
	"kool-dev/kool/services/compose"
	"kool-dev/kool/services/updater"

	"github.com/spf13/cobra"
)

// KoolStartFlags holds the flags for the kool start command
type KoolStartFlags struct {
	Foreground bool
	Rebuild    bool
}

// KoolStart holds handlers and functions for starting containers logic
type KoolStart struct {
	DefaultKoolService
	Flags *KoolStartFlags

	check      checker.Checker
	net        network.Handler
	envStorage environment.EnvStorage
	start      builder.Command

	rebuilder KoolService
}

// KoolRebuild holds handlers for updating the service's images
type KoolRebuild struct {
	DefaultKoolService

	pull, build builder.Command
}

// NewStartCommand initializes new kool start Cobra command
func NewStartCommand(start *KoolStart) (startCmd *cobra.Command) {
	startCmd = &cobra.Command{
		Use:   "start [SERVICE...]",
		Short: "Start service containers defined in docker-compose.yml",
		Long: `Start one or more specified [SERVICE] containers. If no [SERVICE] is provided,
all containers are started. If the containers are already running, they are recreated.`,
		Run: DefaultCommandRunFunction(CheckNewVersion(start, &updater.DefaultUpdater{RootCommand: rootCmd})),

		DisableFlagsInUseLine: true,
	}

	startCmd.Flags().BoolVarP(&start.Flags.Foreground, "foreground", "f", false, "Start containers in foreground mode")
	startCmd.Flags().BoolVarP(&start.Flags.Rebuild, "rebuild", "b", false, "Updates and builds service's images")

	return
}

// NewKoolStart creates a new pointer with default KoolStart service
// dependencies.
func NewKoolStart() *KoolStart {
	defaultKoolService := newDefaultKoolService()
	return &KoolStart{
		*defaultKoolService,
		&KoolStartFlags{false, false},
		checker.NewChecker(defaultKoolService.shell),
		network.NewHandler(defaultKoolService.shell),
		environment.NewEnvStorage(),
		compose.NewDockerCompose("up", "--force-recreate"),
		&KoolRebuild{
			*newDefaultKoolService(),
			compose.NewDockerCompose("pull"),
			compose.NewDockerCompose("build", "--pull"),
		},
	}
}

func AddKoolStart(root *cobra.Command) {
	root.AddCommand(NewStartCommand(NewKoolStart()))
}

// Execute runs the rebuild logic
func (r *KoolRebuild) Execute(args []string) (err error) {
	if err = r.Interactive(r.pull); err != nil {
		return
	}

	err = r.Interactive(r.build)
	return
}

// Execute runs the start logic with incoming arguments
func (s *KoolStart) Execute(args []string) (err error) {
	if s.Flags.Rebuild {
		if err = s.rebuild(); err != nil {
			return
		}
	}

	if !s.Flags.Foreground {
		s.start.AppendArgs("-d")
	}

	if err = s.checkDependencies(); err != nil {
		return
	}

	err = s.Interactive(s.start, args...)
	return
}

func (s *KoolStart) rebuild() (err error) {
	var task = NewKoolTask("Pulling and building service's images", s.rebuilder)

	task.SetInStream(s.InStream())
	task.SetOutStream(s.OutStream())
	task.SetErrStream(s.ErrStream())

	err = task.Run(nil)
	return
}

func (s *KoolStart) checkDependencies() (err error) {
	chErrDocker, chErrNetwork := s.checkDocker(), s.checkNetwork()
	errDocker, errNetwork := <-chErrDocker, <-chErrNetwork

	if errDocker != nil {
		err = errDocker
		return
	}

	if errNetwork != nil {
		err = errNetwork
		return
	}

	return
}

func (s *KoolStart) checkDocker() <-chan error {
	err := make(chan error)

	go func() {
		err <- s.check.Check()
	}()

	return err
}

func (s *KoolStart) checkNetwork() <-chan error {
	err := make(chan error)

	go func() {
		err <- s.net.HandleGlobalNetwork(s.envStorage.Get("KOOL_GLOBAL_NETWORK"))
	}()

	return err
}
