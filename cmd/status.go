package cmd

import (
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/shell"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the status for containers",
	Run:   runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	var dependenciesChecker = checker.NewChecker()

	if err := dependenciesChecker.VerifyDependencies(); err != nil {
		shell.ExecError("", err)
		os.Exit(1)
	}

	handleGlobalNetwork()
	statusDisplayServices()
}

type statusService struct {
	service, state, ports string
	running               bool
}

func statusDisplayServices() {
	out, err := shell.Exec("docker-compose", "ps", "--services")

	if err != nil {
		shell.Warning("No services found.")
		return
	}

	parsedServices := strings.Split(strings.Replace(out, "\r\n", "\n", -1), "\n")
	services := []string{}
	for _, s := range parsedServices {
		if s != "" {
			services = append(services, s)
		}
	}
	if len(services) == 0 {
		shell.Warning("No services found.")
		return
	}

	chStatus := make(chan *statusService, len(services))

	for _, service := range services {
		go func(service string, ch chan *statusService) {
			ss := &statusService{service: service}

			out, err = shell.Exec("docker-compose", "ps", "-q", service)

			if err != nil {
				shell.ExecError(out, err)
				os.Exit(1)
			}

			if out != "" {
				out, err = shell.Exec("docker", "ps", "-a", "--filter", "ID="+out, "--format", "{{.Status}}|{{.Ports}}")

				containerInfo := strings.Split(out, "|")
				ss.running = strings.HasPrefix(containerInfo[0], "Up")
				ss.state = containerInfo[0]
				if len(containerInfo) > 1 {
					ss.ports = containerInfo[1]
				}
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
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Service", "Running", "Ports", "State"})

	for _, s := range status {
		running := "Not running"
		if s.running {
			running = "Running"
		}
		t.AppendRow([]interface{}{s.service, running, s.ports, s.state})
	}

	t.Render()
}
