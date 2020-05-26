package cmd

import (
	"fmt"
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
	handleGlobalNetwork()

	statusDisplayServices()
}

type statusService struct {
	service, state, ports string
	running               bool
}

func statusDisplayServices() {
	out, err := shellExec("docker-compose", "ps", "--services")

	if err != nil {
		fmt.Println("No services found.")
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
		fmt.Println("No services found.")
		return
	}

	status := make([]*statusService, len(services))

	for i, service := range services {
		status[i] = &statusService{service: service}
		out, err = shellExec("docker-compose", "ps", "-q", service)

		if err != nil {
			execError(out, err)
			os.Exit(1)
		}

		if out != "" {
			status[i].running = true
			// it is running
			out, err = shellExec("docker", "ps", "-a", "--filter", "ID="+out, "--format", "{{.Status}}|{{.Ports}}")

			containerInfo := strings.Split(out, "|")
			status[i].state = containerInfo[0]
			if len(containerInfo) > 1 {
				status[i].ports = containerInfo[1]
			}
			containerInfo = nil
		}
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
