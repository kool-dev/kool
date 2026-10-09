package commands

import (
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/network"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/checker"
	"sort"
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
	getProjectsCmd          builder.Command
	getServiceIDCmd         builder.Command
	getProjectServiceIDCmd  builder.Command
	getServiceStatusPortCmd builder.Command

	table shell.TableWriter
}

type statusService struct {
	project, service, state, ports string
	running                        string
	err                            error
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
		builder.NewCommand("docker", "compose", "config", "--services"),
		builder.NewCommand("docker", "ps", "--all"),
		builder.NewCommand("docker", "compose", "ps", "--all", "--quiet"),
		builder.NewCommand("docker", "ps", "--all", "--quiet"),
		builder.NewCommand("docker", "ps", "--all", "--format", "{{.Status}}|{{.Ports}}"),
		shell.NewTableWriter(),
	}
}

// Execute runs the status logic with incoming arguments.
func (s *KoolStatus) Execute(args []string) (err error) {
	if !workspacesEnabled(s.env) {
		return s.executeLegacy()
	}
	var projects []statusProject

	if err = s.checkDependencies(); err != nil {
		return
	}

	if projects, err = s.getProjects(); err != nil {
		return
	} else if len(projects) == 0 {
		s.Shell().Warning("No services found.")
		return
	}
	serviceCount := 0
	for _, project := range projects {
		serviceCount += len(project.services)
	}
	if serviceCount == 0 {
		s.Shell().Warning("No services found.")
		return
	}

	chStatus := make(chan *statusService, serviceCount)

	s.table.SetWriter(s.Shell().OutStream())
	s.table.AppendHeader("Project", "Service", "Running", "Ports", "State")

	go func() {
		var wg sync.WaitGroup

		defer close(chStatus)

		for _, project := range projects {
			for _, service := range project.services {
				wg.Add(1)
				go s.fetchServiceInfo(project.name, service, chStatus, &wg)
			}
		}

		wg.Wait()
	}()

	var statuses []*statusService
	for ss := range chStatus {
		if ss.err != nil {
			err = ss.err
			return
		}
		statuses = append(statuses, ss)
	}
	sort.Slice(statuses, func(i, j int) bool {
		if statuses[i].project == statuses[j].project {
			return statuses[i].service < statuses[j].service
		}
		return statuses[i].project < statuses[j].project
	})
	for _, ss := range statuses {
		s.table.AppendRow(ss.project, ss.service, ss.running, ss.ports, ss.state)
	}
	s.table.Render()
	return
}

func (s *KoolStatus) executeLegacy() (err error) {
	if err = s.checkDependencies(); err != nil {
		return
	}
	services, err := s.getServices()
	if err != nil {
		return err
	}
	if len(services) == 0 {
		s.Shell().Warning("No services found.")
		return nil
	}

	chStatus := make(chan *statusService, len(services))
	s.table.SetWriter(s.Shell().OutStream())
	s.table.AppendHeader("Service", "Running", "Ports", "State")
	go func() {
		var wg sync.WaitGroup
		defer close(chStatus)
		for _, service := range services {
			wg.Add(1)
			go s.fetchLegacyServiceInfo(service, chStatus, &wg)
		}
		wg.Wait()
	}()
	for status := range chStatus {
		if status.err != nil {
			return status.err
		}
		s.table.AppendRow(status.service, status.running, status.ports, status.state)
	}
	s.table.SortBy(1)
	s.table.Render()
	return nil
}

type statusProject struct {
	name     string
	services []string
}

func (s *KoolStatus) getProjects() ([]statusProject, error) {
	services, err := s.getServices()
	if err != nil {
		return nil, err
	}
	mainProject := sourceProject(s.env)
	projects := []statusProject{{name: mainProject, services: services}}
	workspaceProjectNames := []string{}
	if isWorkspace(s.env) {
		workspaceProjectNames = append(workspaceProjectNames, currentProject(s.env))
	} else if workspaceProjectNames, err = activeWorkspaceProjects(s.Shell(), s.getProjectsCmd, s.env); err != nil {
		return nil, err
	}
	workspaceProjectServices := configuredWorkspaceServices(s.env)
	for _, project := range workspaceProjectNames {
		if project != "" && project != mainProject {
			projects = append(projects, statusProject{name: project, services: workspaceProjectServices})
		}
	}
	return projects, nil
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

	parsedServices := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")
	for _, s := range parsedServices {
		if s != "" {
			services = append(services, s)
		}
	}
	return
}

func (s *KoolStatus) fetchServiceInfo(project, service string, chStatus chan *statusService, wg *sync.WaitGroup) {
	var isRunning bool

	defer wg.Done()

	ss := &statusService{project: project, service: service, running: "Not running"}
	isRunning, ss.state, ss.ports, ss.err = s.getServiceInfo(project, service)
	if isRunning {
		ss.running = "Running"
	}

	chStatus <- ss
}

func (s *KoolStatus) fetchLegacyServiceInfo(service string, chStatus chan *statusService, wg *sync.WaitGroup) {
	defer wg.Done()
	status := &statusService{service: service, running: "Not running"}
	var serviceID string
	if serviceID, status.err = s.Shell().Exec(s.getServiceIDCmd, service); status.err == nil && serviceID != "" {
		status.state, status.ports = s.getStatusPort(serviceID)
		if strings.HasPrefix(status.state, "Up") {
			status.running = "Running"
		}
	}
	chStatus <- status
}

func (s *KoolStatus) getServiceInfo(project, service string) (isRunning bool, status, port string, err error) {
	var serviceID string
	if serviceID, err = s.Shell().Exec(s.getProjectServiceIDCmd, "--filter", "label=com.docker.compose.project="+project, "--filter", "label=com.docker.compose.service="+service, "--filter", "label=com.docker.compose.oneoff=False"); err == nil && serviceID != "" {
		serviceID = strings.Fields(serviceID)[0]
		status, port = s.getStatusPort(serviceID)
		if strings.HasPrefix(status, "Up") {
			isRunning = true
		}
	}
	return
}

func (s *KoolStatus) getCurrentServiceInfo(service string) (bool, string, string, error) {
	if workspacesEnabled(s.env) {
		return s.getServiceInfo(currentProject(s.env), service)
	}
	serviceID, err := s.Shell().Exec(s.getServiceIDCmd, service)
	if err != nil || serviceID == "" {
		return false, "", "", err
	}
	status, port := s.getStatusPort(serviceID)
	return strings.HasPrefix(status, "Up"), status, port, nil
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
