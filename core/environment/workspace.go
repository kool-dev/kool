package environment

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/compose-spec/compose-go/template"
	"gopkg.in/yaml.v2"
)

var workspaceOverrideFile string

// WorkspaceSourceLabel identifies containers owned by a source project's workspaces.
const WorkspaceSourceLabel = "dev.kool.workspace.source"

func initWorkspaceContext(envStorage EnvStorage, workDir, source, workspace, provider string, active, uniqueHost bool) {
	envStorage.Set("KOOL_WORKSPACE_PROVIDER", provider)
	envStorage.Set("KOOL_WORKSPACE_SOURCE", source)
	if envStorage.Get("KOOL_NAME") == "" {
		envStorage.Set("KOOL_NAME", filepath.Base(source))
	}

	sourceProject := envStorage.Get("COMPOSE_PROJECT_NAME")
	if sourceProject == "" {
		sourceProject = composeSourceProject(envStorage, workDir)
		if sourceProject == "" {
			sourceProject = composeProjectName(filepath.Base(source))
		}
	}
	envStorage.Set("KOOL_WORKSPACE_SOURCE_PROJECT", sourceProject)
	if !active {
		return
	}

	workspaceName := filepath.Base(workspace)
	workspaceProject := sourceProject + "-workspace-" + composeProjectName(workspaceIdentity(workspace))
	workspaceHost := workspaceName
	if uniqueHost {
		workspaceHost = workspaceIdentity(workspace)
	}
	envStorage.Set("KOOL_WORKSPACE", "true")
	envStorage.Set("KOOL_WORKSPACE_NAME", workspaceHostName(workspaceHost))
	envStorage.Set("KOOL_WORKSPACE_PATH", workspace)
	envStorage.Set("KOOL_WORKSPACE_PROJECT", workspaceProject)
	envStorage.Set("COMPOSE_PROJECT_NAME", workspaceProject)
	initWorkspaceCompose(envStorage, workDir)
}

func workspaceIdentity(workspace string) string {
	canonical, err := filepath.EvalSymlinks(workspace)
	if err != nil {
		canonical, err = filepath.Abs(workspace)
		if err != nil {
			canonical = workspace
		}
	}
	digest := sha256.Sum256([]byte(canonical))
	suffix := fmt.Sprintf("-%x", digest[:4])
	name := filepath.Base(workspace)
	if len(name) > 63-len(suffix) {
		name = strings.TrimRight(name[:63-len(suffix)], "-_")
	}
	return name + suffix
}

func workspaceHostName(name string) string {
	name = strings.ToLower(name)
	var host strings.Builder
	separator := false
	for _, character := range name {
		if character >= 'a' && character <= 'z' || character >= '0' && character <= '9' {
			host.WriteRune(character)
			separator = false
		} else if host.Len() > 0 && !separator {
			host.WriteByte('-')
			separator = true
		}
	}
	result := strings.Trim(host.String(), "-")
	if result == "" {
		return "workspace"
	}
	if len(result) > 63 {
		result = strings.TrimRight(result[:63], "-")
	}
	return result
}

type workspaceComposeConfig struct {
	Name string `yaml:"name"`
}

func initWorkspaceCompose(envStorage EnvStorage, workDir string) []string {
	koolConfig := loadKoolConfig(workDir)
	if koolConfig == nil || len(koolConfig.Workspaces) == 0 {
		return nil
	}
	unique := make(map[string]bool, len(koolConfig.Workspaces))
	var selected []string
	for _, service := range koolConfig.Workspaces {
		if service != "" && !unique[service] {
			selected = append(selected, service)
			unique[service] = true
		}
	}
	sort.Strings(selected)
	if len(selected) == 0 {
		return nil
	}

	files := composeFiles(envStorage, workDir)
	if len(files) == 0 {
		return nil
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil
		}

		config := workspaceComposeConfig{}
		if yaml.Unmarshal(content, &config) != nil {
			return nil
		}
	}

	CleanupWorkspace()
	override, err := os.CreateTemp("", "kool-workspace-*.yml")
	if err != nil {
		return nil
	}

	var content strings.Builder
	content.WriteString("services:\n")
	for _, name := range selected {
		fmt.Fprintf(&content, "  %q:\n", name)
		content.WriteString("    ports: !reset []\n")
		content.WriteString("    container_name: !reset null\n")
		content.WriteString("    labels:\n")
		fmt.Fprintf(&content, "      %q: %q\n", WorkspaceSourceLabel, envStorage.Get("KOOL_WORKSPACE_SOURCE_PROJECT"))
	}

	if _, err = override.WriteString(content.String()); err != nil {
		_ = override.Close()
		_ = os.Remove(override.Name())
		return nil
	}
	if err = override.Close(); err != nil {
		_ = os.Remove(override.Name())
		return nil
	}
	workspaceOverrideFile = override.Name()

	files = append(files, workspaceOverrideFile)
	separator := envStorage.Get("COMPOSE_PATH_SEPARATOR")
	if separator == "" {
		separator = string(os.PathListSeparator)
	}
	envStorage.Set("COMPOSE_FILE", strings.Join(files, separator))
	envStorage.Set("KOOL_WORKSPACE_SERVICES", strings.Join(selected, ","))
	return selected
}

// CleanupWorkspace removes the generated Compose override for this process.
func CleanupWorkspace() {
	if workspaceOverrideFile != "" {
		_ = os.Remove(workspaceOverrideFile)
		workspaceOverrideFile = ""
	}
}

// IsolateWorkspace gives a recursive command its own generated Compose override.
func IsolateWorkspace() func() {
	parentOverride := workspaceOverrideFile
	workspaceOverrideFile = ""
	return func() {
		CleanupWorkspace()
		workspaceOverrideFile = parentOverride
	}
}

func composeSourceProject(envStorage EnvStorage, workDir string) string {
	project := ""
	for _, file := range composeFiles(envStorage, workDir) {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		config := workspaceComposeConfig{}
		if yaml.Unmarshal(content, &config) == nil && config.Name != "" {
			resolved, err := template.Substitute(config.Name, os.LookupEnv)
			if err == nil {
				project = composeProjectName(resolved)
			}
		}
	}
	return project
}

func composeFiles(envStorage EnvStorage, workDir string) []string {
	if configured := envStorage.Get("COMPOSE_FILE"); configured != "" {
		separator := envStorage.Get("COMPOSE_PATH_SEPARATOR")
		if separator == "" {
			separator = string(os.PathListSeparator)
		}

		files := strings.Split(configured, separator)
		for index, file := range files {
			if !filepath.IsAbs(file) {
				files[index] = filepath.Join(workDir, file)
			}
		}
		return files
	}

	for _, name := range []string{"compose.yaml", "compose.yml", "docker-compose.yml", "docker-compose.yaml"} {
		file := filepath.Join(workDir, name)
		if _, err := os.Stat(file); err == nil {
			return []string{file}
		}
	}
	return nil
}
