package cmd

import (
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/network"
	"kool-dev/kool/cmd/shell"
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// DefaultStatusCmd holds data for status command
type DefaultStatusCmd struct{}

// StatusCmd holds logic for status command
type StatusCmd interface {
	CheckDependencies() error
	GetServices() ([]string, error)
	GetServiceID(service string) (string, error)
	GetStatusPort(serviceID string) (string, string)
}

type statusService struct {
	service, state, ports string
	running               bool
}

// NewStatusCommand Initialize new kool status command
func NewStatusCommand(statusCmd StatusCmd) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Shows the status for containers",
		Run: func(cmd *cobra.Command, args []string) {
			if err := statusCmd.CheckDependencies(); err != nil {
				shell.FexecError(cmd.OutOrStdout(), "", err)
				os.Exit(1)
			}

			statusDisplayServices(statusCmd, cmd)
		},
	}
}

// CheckDependencies check kool dependencies
func (s *DefaultStatusCmd) CheckDependencies() (err error) {
	var dependenciesChecker = checker.NewChecker()

	if err = dependenciesChecker.VerifyDependencies(); err != nil {
		return
	}

	var globalNetworkHandler = network.NewHandler()

	if err = globalNetworkHandler.HandleGlobalNetwork(); err != nil {
		return
	}

	return
}

// GetServices get docker services
func (s *DefaultStatusCmd) GetServices() (services []string, err error) {
	var output string

	cmd := builder.NewCommand("docker-compose", "ps", "--services")

	if output, err = cmd.Exec(); err != nil {
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

// GetServiceID get docker service ID
func (s *DefaultStatusCmd) GetServiceID(service string) (serviceID string, err error) {
	cmd := builder.NewCommand("docker-compose", "ps", "-q", service)
	serviceID, err = cmd.Exec()
	return
}

// GetStatusPort get docker service port and status
func (s *DefaultStatusCmd) GetStatusPort(serviceID string) (status string, port string) {
	var output string

	cmd := builder.NewCommand("docker", "ps", "-a", "--filter", "ID="+serviceID, "--format", "{{.Status}}|{{.Ports}}")

	if output, _ = cmd.Exec(); output == "" {
		return
	}

	containerInfo := strings.Split(output, "|")

	status = containerInfo[0]

	if len(containerInfo) > 1 {
		port = containerInfo[1]
	}

	return
}

func init() {
	rootCmd.AddCommand(NewStatusCommand(&DefaultStatusCmd{}))
}

func statusDisplayServices(statusCmd StatusCmd, cobraCmd *cobra.Command) {
	services, err := statusCmd.GetServices()

	if err != nil {
		shell.Fwarning(cobraCmd.OutOrStdout(), "No services found.")
		return
	}

	if len(services) == 0 {
		shell.Fwarning(cobraCmd.OutOrStdout(), "No services found.")
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

			if serviceID, err = statusCmd.GetServiceID(service); err != nil {
				shell.FexecError(cobraCmd.OutOrStdout(), serviceID, err)
				os.Exit(1)
			}

			if serviceID != "" {
				status, port = statusCmd.GetStatusPort(serviceID)

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
		if i == l-1 {
			close(chStatus)
			break
		}
		i++
	}

	t := table.NewWriter()
	t.SetOutputMirror(cobraCmd.OutOrStdout())
	t.AppendHeader(table.Row{"Service", "Running", "Ports", "State"})

	sort.SliceStable(status, func(i, j int) bool {
		return status[i].service < status[j].service
	})

	for _, s := range status {
		running := "Not running"
		if s.running {
			running = "Running"
		}
		t.AppendRow([]interface{}{s.service, running, s.ports, s.state})
	}

	t.Render()
}
