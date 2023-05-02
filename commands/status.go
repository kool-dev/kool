package commands

import (
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/network"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/checker"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// KoolStatus holds handlers and functions to implement the status command logic
type KoolStatus struct {
	DefaultKoolService

	check checker.Checker
	net   network.Handler
	env   environment.EnvStorage

	getServicesCmd          builder.Command
	getServiceIDCmd         builder.Command
	getServiceStatusPortCmd builder.Command

	table shell.TableWriter
}

type statusService struct {
	service, state, ports string
	running               string
	err                   error
}

func AddKoolStatus(root *cobra.Command) {
	var (
		status    = NewKoolStatus()
		statusCmd = NewStatusCommand(status)
	)

	root.AddCommand(statusCmd)
}

// NewKoolStatus creates a new handler for status logic
func NewKoolStatus() *KoolStatus {
	defaultKoolService := newDefaultKoolService()
	return &KoolStatus{
		*defaultKoolService,
		checker.NewChecker(defaultKoolService.shell),
		network.NewHandler(defaultKoolService.shell),
		environment.NewEnvStorage(),
		builder.NewCommand("docker", "compose", "ps", "--all", "--services"),
		builder.NewCommand("docker", "compose", "ps", "--all", "--quiet"),
		builder.NewCommand("docker", "ps", "--all", "--format", "{{.Status}}|{{.Ports}}"),
		shell.NewTableWriter(),
	}
}

// Execute runs the status logic with incoming arguments.
func (s *KoolStatus) Execute(args []string) (err error) {
	var services []string

	if err = s.checkDependencies(); err != nil {
		return
	}

	if services, err = s.getServices(); err != nil || len(services) == 0 {
		s.Shell().Warning("No services found.")
		return
	}

	chStatus := make(chan *statusService, len(services))

	s.table.SetWriter(s.Shell().OutStream())
	s.table.AppendHeader("Service", "Running", "Ports", "State")

	go func() {
		var wg sync.WaitGroup

		defer close(chStatus)

		for _, service := range services {
			wg.Add(1)
			go s.fetchServiceInfo(service, chStatus, &wg)
		}

		wg.Wait()
	}()

	for ss := range chStatus {
		if ss.err != nil {
			err = ss.err
			return
		}

		s.table.AppendRow(ss.service, ss.running, ss.ports, ss.state)
	}

	s.table.SortBy(1)
	s.table.Render()
	return
}

func (s *KoolStatus) checkDependencies() (err error) {
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

func (s *KoolStatus) checkDocker() <-chan error {
	err := make(chan error)

	go func() {
		err <- s.check.Check()
	}()

	return err
}

func (s *KoolStatus) checkNetwork() <-chan error {
	err := make(chan error)

	go func() {
		err <- s.net.HandleGlobalNetwork(s.env.Get("KOOL_GLOBAL_NETWORK"))
	}()

	return err
}

func (s *KoolStatus) getServices() (services []string, err error) {
	var output string

	if output, err = s.Shell().Exec(s.getServicesCmd); err != nil {
		return
	}

	parsedServices := strings.Split(strings.Replace(output, "\r\n", "\n", -1), "\n")
	for _, s := range parsedServices {
		if s != "" {
			services = append(services, s)
		}
	}

	return
}

func (s *KoolStatus) fetchServiceInfo(service string, chStatus chan *statusService, wg *sync.WaitGroup) {
	var isRunning bool

	defer wg.Done()

	ss := &statusService{service: service, running: "Not running"}
	isRunning, ss.state, ss.ports, ss.err = s.getServiceInfo(service)
	if isRunning {
		ss.running = "Running"
	}

	chStatus <- ss
}

func (s *KoolStatus) getServiceInfo(service string) (isRunning bool, status, port string, err error) {
	var serviceID string
	if serviceID, err = s.Shell().Exec(s.getServiceIDCmd, service); err == nil && serviceID != "" {
		status, port = s.getStatusPort(serviceID)
		if strings.HasPrefix(status, "Up") {
			isRunning = true
		}
	}
	return
}

func (s *KoolStatus) getStatusPort(serviceID string) (status string, port string) {
	var output string

	if output, _ = s.Shell().Exec(s.getServiceStatusPortCmd, "--filter", "ID="+serviceID); output == "" {
		return
	}

	containerInfo := strings.Split(output, "|")

	status = containerInfo[0]

	if len(containerInfo) > 1 {
		port = containerInfo[1]
	}

	return
}

// NewStatusCommand Initialize new kool status command
func NewStatusCommand(status *KoolStatus) *cobra.Command {
	var statusTask = NewKoolTask("Fetching services status", status)

	statusTask.SetFrameOutput(false)

	return &cobra.Command{
		Use:   "status",
		Short: "Show the status of all service containers",
		RunE:  LongTaskCommandRunFunction(statusTask),

		DisableFlagsInUseLine: true,
	}
}
