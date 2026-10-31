package commands

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/parser"
	"kool-dev/kool/core/shell"
	"sort"
	"strings"
)

func isWorkspace(env environment.EnvStorage) bool {
	return env.IsTrue("KOOL_WORKSPACE")
}

func workspacesEnabled(env environment.EnvStorage) bool {
	return env.IsTrue("KOOL_WORKSPACES_ENABLED")
}

func proxyEnabled(env environment.EnvStorage) bool {
	return env.IsTrue("KOOL_PROXY_ENABLED")
}

func workspaceServices(env environment.EnvStorage) []string {
	configured := env.Get("KOOL_WORKSPACE_SERVICES")
	if configured == "" {
		return nil
	}
	return strings.Split(configured, ",")
}

func configuredWorkspaceServices(env environment.EnvStorage) []string {
	if services := workspaceServices(env); len(services) > 0 {
		return services
	}
	config, err := parser.LoadKoolYaml(env.Get("PWD"))
	if err != nil {
		return nil
	}
	services := append([]string(nil), config.Workspaces...)
	sort.Strings(services)
	return services
}

func sourceProject(env environment.EnvStorage) string {
	if project := env.Get("KOOL_WORKSPACE_SOURCE_PROJECT"); project != "" {
		return project
	}
	if project := env.Get("COMPOSE_PROJECT_NAME"); project != "" && !isWorkspace(env) {
		return project
	}
	return env.Get("KOOL_NAME")
}

func currentProject(env environment.EnvStorage) string {
	if isWorkspace(env) {
		return env.Get("KOOL_WORKSPACE_PROJECT")
	}
	return sourceProject(env)
}

func activeWorkspaceProjects(sh shell.Shell, command builder.Command, env environment.EnvStorage) ([]string, error) {
	output, err := sh.Exec(command,
		"--filter", "label="+environment.WorkspaceSourceLabel+"="+sourceProject(env),
		"--format", `{{.Label "com.docker.compose.project"}}`,
	)
	if err != nil {
		return nil, err
	}
	prefix := sourceProject(env) + "-workspace-"
	seen := make(map[string]bool)
	var projects []string
	for _, project := range strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n") {
		project = strings.TrimSpace(project)
		if strings.HasPrefix(project, prefix) && !seen[project] {
			projects = append(projects, project)
			seen[project] = true
		}
	}
	sort.Strings(projects)
	return projects, nil
}

func selectWorkspaceServices(env environment.EnvStorage, requested []string) ([]string, error) {
	if !isWorkspace(env) {
		return requested, nil
	}

	configured := workspaceServices(env)
	if len(configured) == 0 {
		return nil, fmt.Errorf("no services are configured under workspaces in kool.yml")
	}
	if len(requested) == 0 {
		return configured, nil
	}

	allowed := make(map[string]bool, len(configured))
	for _, service := range configured {
		allowed[service] = true
	}
	for _, service := range requested {
		if !allowed[service] {
			return nil, fmt.Errorf("service %s is not enabled for workspaces", service)
		}
	}
	return requested, nil
}
