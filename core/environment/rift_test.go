package environment

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitRiftWorkspace(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "My App")
	workspace := filepath.Join(root, ".rifts", "My App", "task-a")
	workDir := filepath.Join(workspace, "packages", "api")

	for _, directory := range []string{source, workDir} {
		if err := os.MkdirAll(directory, 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(workspace, ".rift"), []byte("workspace"), 0644); err != nil {
		t.Fatal(err)
	}

	originalRiftAncestors := riftAncestors
	riftAncestors = func(string) ([]byte, error) {
		return []byte(source + "\n"), nil
	}
	defer func() { riftAncestors = originalRiftAncestors }()

	env := NewFakeEnvStorage()
	initRift(env, workDir)

	if !env.IsTrue("KOOL_WORKSPACE") {
		t.Error("expected created Rift to be a workspace")
	}
	if got := env.Get("KOOL_WORKSPACE_PROVIDER"); got != "rift" {
		t.Errorf("expected Rift provider, got %q", got)
	}
	if got := env.Get("KOOL_WORKSPACE_SOURCE"); got != source {
		t.Errorf("expected KOOL_WORKSPACE_SOURCE %q, got %q", source, got)
	}
	if got := env.Get("KOOL_WORKSPACE_PATH"); got != workspace {
		t.Errorf("expected KOOL_WORKSPACE_PATH %q, got %q", workspace, got)
	}
	expectedProject := "myapp-workspace-" + composeProjectName(workspaceIdentity(workspace))
	if got := env.Get("COMPOSE_PROJECT_NAME"); got != expectedProject {
		t.Errorf("expected COMPOSE_PROJECT_NAME %q, got %q", expectedProject, got)
	}
}

func TestWorkspaceHostNameIsDNSLabel(t *testing.T) {
	for input, expected := range map[string]string{
		"Feature Login": "feature-login",
		"TASK_api":      "task-api",
		"---":           "workspace",
	} {
		if actual := workspaceHostName(input); actual != expected {
			t.Errorf("expected workspace host %q for %q, got %q", expected, input, actual)
		}
	}
	if actual := workspaceHostName(strings.Repeat("a", 70)); len(actual) != 63 {
		t.Errorf("expected workspace host to be limited to 63 characters, got %d", len(actual))
	}
}

func TestInitRiftUsesConfiguredSourceProjectName(t *testing.T) {
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, ".rift"), []byte("workspace"), 0644); err != nil {
		t.Fatal(err)
	}

	originalRiftAncestors := riftAncestors
	riftAncestors = func(string) ([]byte, error) {
		return []byte("/projects/app\n"), nil
	}
	defer func() { riftAncestors = originalRiftAncestors }()

	env := NewFakeEnvStorage()
	env.Set("COMPOSE_PROJECT_NAME", "custom-project")
	initRift(env, workspace)

	expected := "custom-project-workspace-" + composeProjectName(workspaceIdentity(workspace))
	if got := env.Get("COMPOSE_PROJECT_NAME"); got != expected {
		t.Errorf("expected COMPOSE_PROJECT_NAME %q, got %q", expected, got)
	}
}

func TestInitRiftWorkspaceComposeProject(t *testing.T) {
	t.Cleanup(CleanupWorkspace)
	root := t.TempDir()
	source := filepath.Join(root, "app")
	workspace := filepath.Join(root, ".rifts", "app", "task-a")
	if err := os.MkdirAll(workspace, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, ".rift"), []byte("workspace"), 0644); err != nil {
		t.Fatal(err)
	}
	compose := `name: custom-app
services:
  database:
    image: mysql
  node:
    image: node
  app:
    image: php
    ports:
      - "80:80"
`
	if err := os.WriteFile(filepath.Join(workspace, "compose.yml"), []byte(compose), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "kool.yml"), []byte("workspaces: [app, node]\n"), 0644); err != nil {
		t.Fatal(err)
	}

	originalRiftAncestors := riftAncestors
	riftAncestors = func(string) ([]byte, error) { return []byte(source + "\n"), nil }
	defer func() { riftAncestors = originalRiftAncestors }()

	env := NewFakeEnvStorage()
	initRift(env, workspace)

	if !env.IsTrue("KOOL_WORKSPACE") {
		t.Error("expected KOOL_WORKSPACE to be enabled")
	}
	expectedProject := "custom-app-workspace-" + composeProjectName(workspaceIdentity(workspace))
	if got := env.Get("COMPOSE_PROJECT_NAME"); got != expectedProject {
		t.Errorf("expected workspace Compose project %q, got %q", expectedProject, got)
	}
	if got := env.Get("KOOL_WORKSPACE_SERVICES"); got != "app,node" {
		t.Errorf("expected sorted workspace services, got %q", got)
	}
	files := strings.Split(env.Get("COMPOSE_FILE"), string(os.PathListSeparator))
	if len(files) != 2 {
		t.Fatalf("expected Compose file and generated override, got %v", files)
	}
	override, err := os.ReadFile(files[1])
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"ports: !reset []", "container_name: !reset null"} {
		if !strings.Contains(string(override), expected) {
			t.Errorf("expected generated override to contain %q, got:\n%s", expected, override)
		}
	}
}

func TestInitRiftPreservesComposeFileProjectName(t *testing.T) {
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, ".rift"), []byte("workspace"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "compose.yml"), []byte("name: custom-project\nservices: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	originalRiftAncestors := riftAncestors
	riftAncestors = func(string) ([]byte, error) {
		return []byte("/projects/app\n"), nil
	}
	defer func() { riftAncestors = originalRiftAncestors }()

	env := NewFakeEnvStorage()
	initRift(env, workspace)

	expected := "custom-project-workspace-" + composeProjectName(workspaceIdentity(workspace))
	if got := env.Get("COMPOSE_PROJECT_NAME"); got != expected {
		t.Errorf("expected top-level Compose project name %q, got %q", expected, got)
	}
}

func TestInitRiftOriginalWorkspaceSetsSourceContext(t *testing.T) {
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, ".rift"), []byte("workspace"), 0644); err != nil {
		t.Fatal(err)
	}

	originalRiftAncestors := riftAncestors
	riftAncestors = func(string) ([]byte, error) { return nil, nil }
	defer func() { riftAncestors = originalRiftAncestors }()

	env := NewFakeEnvStorage()
	initRift(env, workspace)

	if env.IsTrue("KOOL_WORKSPACE") {
		t.Error("expected original workspace not to be a created workspace")
	}
	if got := env.Get("KOOL_WORKSPACE_SOURCE"); got != workspace {
		t.Errorf("expected KOOL_WORKSPACE_SOURCE %q, got %q", workspace, got)
	}
}
