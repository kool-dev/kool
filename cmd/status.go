package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/network"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"sort"
	"strings"

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
	running               bool
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
	if err = s.check.Check(); err != nil {
		return
	}

	if err = s.net.HandleGlobalNetwork(s.envStorage.Get("KOOL_GLOBAL_NETWORK")); err != nil {
		return
	}

	services, err := s.getServices()

	if err != nil {
		s.Warning("No services found.")
		return
	}

	if len(services) == 0 {
		s.Warning("No services found.")
		return
	}

	chStatus := make(chan *statusService, len(services))

	for _, service := range services {
		go func(service string, ch chan *statusService) {
			var (
				status, port, serviceID string
				err                     error
			)

			ss := &statusService{service: service}

			if serviceID, err = s.Exec(s.getServiceIDRunner, service); err != nil {
				ss.err = err
			} else if serviceID != "" {
				status, port = s.getStatusPort(serviceID)

				ss.running = strings.HasPrefix(status, "Up")
				ss.state = status
				ss.ports = port
			}

			ch <- ss
		}(service, chStatus)
	}

	var i, l int = 0, len(services)
	status := make([]*statusService, l)
	for ss := range chStatus {
		status[i] = ss

		if status[i].err != nil {
			err = status[i].err
			return
		}

		if i == l-1 {
			close(chStatus)
			break
		}
		i++
	}

	s.table.SetWriter(s.OutStream())
	s.table.AppendHeader("Service", "Running", "Ports", "State")

	sort.SliceStable(status, func(i, j int) bool {
		return status[i].service < status[j].service
	})

	for _, st := range status {
		running := "Not running"
		if st.running {
			running = "Running"
		}
		s.table.AppendRow(st.service, running, st.ports, st.state)
	}

	s.table.Render()
	return
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
