package commands

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/parser"
	"kool-dev/kool/core/shell"
	"path/filepath"
	"sort"
	"strings"
)

func isWorkspace(env environment.EnvStorage) bool {
	return env.IsTrue("KOOL_WORKSPACE")
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
	for _, name := range []string{"kool.yml", "kool.yaml"} {
		config, err := parser.ParseKoolYaml(filepath.Join(env.Get("PWD"), name))
		if err == nil {
			services := append([]string(nil), config.Workspaces...)
			sort.Strings(services)
			return services
		}
	}
	return nil
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
