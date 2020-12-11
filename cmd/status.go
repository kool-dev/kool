package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/network"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// KoolStatus holds handlers and functions to implement the status command logic
type KoolStatus struct {
	DefaultKoolService

	check      checker.Checker
	net        network.Handler
	envStorage environment.EnvStorage

	getServicesRunner          builder.Command
	getServiceIDRunner         builder.Command
	getServiceStatusPortRunner builder.Command

	table shell.TableWriter
}

type statusService struct {
	service, state, ports string
	running               string
	err                   error
}

func init() {
	var (
		status    = NewKoolStatus()
		statusCmd = NewStatusCommand(status)
	)

	rootCmd.AddCommand(statusCmd)
}

// NewKoolStatus creates a new handler for status logic
func NewKoolStatus() *KoolStatus {
	defaultKoolService := newDefaultKoolService()
	return &KoolStatus{
		*defaultKoolService,
		checker.NewChecker(defaultKoolService.shell),
		network.NewHandler(defaultKoolService.shell),
		environment.NewEnvStorage(),
		builder.NewCommand("docker-compose", "ps", "--services"),
		builder.NewCommand("docker-compose", "ps", "-q"),
		builder.NewCommand("docker", "ps", "-a", "--format", "{{.Status}}|{{.Ports}}"),
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
		s.Warning("No services found.")
		return
	}

	chStatus := make(chan *statusService, len(services))

	s.table.SetWriter(s.OutStream())
	s.table.AppendHeader("Service", "Running", "Ports", "State")

	go func() {
		var wg sync.WaitGroup

		defer close(chStatus)

		for _, service := range services {
			wg.Add(1)
			go s.getServiceInfo(service, chStatus, &wg)
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
		err <- s.net.HandleGlobalNetwork(s.envStorage.Get("KOOL_GLOBAL_NETWORK"))
	}()

	return err
}

func (s *KoolStatus) getServices() (services []string, err error) {
	var output string

	if output, err = s.Exec(s.getServicesRunner); err != nil {
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

func (s *KoolStatus) getServiceInfo(service string, chStatus chan *statusService, wg *sync.WaitGroup) {
	var (
		status, port, serviceID string
		err                     error
	)

	defer wg.Done()

	ss := &statusService{service: service}

	if serviceID, err = s.Exec(s.getServiceIDRunner, service); err != nil {
		ss.err = err
	} else if serviceID != "" {
		status, port = s.getStatusPort(serviceID)

		ss.running = "Not running"
		if strings.HasPrefix(status, "Up") {
			ss.running = "Running"
		}
		ss.state = status
		ss.ports = port
	}

	chStatus <- ss
}

func (s *KoolStatus) getStatusPort(serviceID string) (status string, port string) {
	var output string

	if output, _ = s.Exec(s.getServiceStatusPortRunner, "--filter", "ID="+serviceID); output == "" {
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
	return &cobra.Command{
		Use:   "status",
		Short: "Shows the status for containers",
		Run:   DefaultCommandRunFunction(status),
	}
}
