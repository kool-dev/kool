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

// DefaultStatusCmd holds interfaces for status command logic
type DefaultStatusCmd struct {
	DependenciesChecker        checker.Checker
	NetworkHandler             network.Handler
	GetServicesRunner          builder.Runner
	GetServiceIDRunner         builder.Runner
	GetServiceStatusPortRunner builder.Runner
}

type statusService struct {
	service, state, ports string
	running               bool
}

// NewStatusCommand Initialize new kool status command
func NewStatusCommand(statusCmd *DefaultStatusCmd) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Shows the status for containers",
		Run: func(cmd *cobra.Command, args []string) {
			if err := checkDependencies(statusCmd); err != nil {
				shell.FexecError(cmd.OutOrStdout(), "", err)
				os.Exit(1)
			}

			statusDisplayServices(statusCmd, cmd)
		},
	}
}

func (s *DefaultStatusCmd) getServices() (services []string, err error) {
	var output string

	if output, err = s.GetServicesRunner.Exec(); err != nil {
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

func (s *DefaultStatusCmd) getStatusPort(serviceID string) (status string, port string) {
	var output string

	if output, _ = s.GetServiceStatusPortRunner.Exec("--filter", "ID="+serviceID); output == "" {
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
	defaultStatusCmd := &DefaultStatusCmd{
		checker.NewChecker(),
		network.NewHandler(),
		builder.NewCommand("docker-compose", "ps", "--services"),
		builder.NewCommand("docker-compose", "ps", "-q"),
		builder.NewCommand("docker", "ps", "-a", "--format", "{{.Status}}|{{.Ports}}"),
	}
	rootCmd.AddCommand(NewStatusCommand(defaultStatusCmd))
}

func checkDependencies(statusCmd *DefaultStatusCmd) (err error) {
	if err = statusCmd.DependenciesChecker.VerifyDependencies(); err != nil {
		return
	}

	if err = statusCmd.NetworkHandler.HandleGlobalNetwork(os.Getenv("KOOL_GLOBAL_NETWORK")); err != nil {
		return
	}

	return
}

func statusDisplayServices(statusCmd *DefaultStatusCmd, cobraCmd *cobra.Command) {
	services, err := statusCmd.getServices()

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

			if serviceID, err = statusCmd.GetServiceIDRunner.Exec(service); err != nil {
				shell.FexecError(cobraCmd.OutOrStdout(), serviceID, err)
				os.Exit(1)
			}

			if serviceID != "" {
				status, port = statusCmd.getStatusPort(serviceID)

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
