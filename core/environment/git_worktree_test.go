package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitGitWorktree(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "My App")
	workspace := filepath.Join(root, "trees", "task-a")
	workDir := filepath.Join(workspace, "packages", "api")
	gitDir := filepath.Join(source, ".git", "worktrees", "task-a")
	commonDir := filepath.Join(source, ".git")

	originalGitWorktreeOutput := gitWorktreeOutput
	gitWorktreeOutput = fakeGitWorktreeOutput(map[string]string{
		"rev-parse --show-toplevel":    workspace,
		"rev-parse --git-dir":          gitDir,
		"rev-parse --git-common-dir":   commonDir,
		"worktree list --porcelain -z": "worktree " + source + "\x00HEAD abc\x00\x00worktree " + workspace + "\x00HEAD def\x00",
	})
	defer func() { gitWorktreeOutput = originalGitWorktreeOutput }()

	env := NewFakeEnvStorage()
	initGitWorktree(env, workDir)

	if !env.IsTrue("KOOL_WORKSPACE") {
		t.Error("expected linked Git worktree to be a workspace")
	}
	if got := env.Get("KOOL_WORKSPACE_PROVIDER"); got != "worktree" {
		t.Errorf("expected worktree provider, got %q", got)
	}
	if got := env.Get("KOOL_WORKSPACE_SOURCE"); got != source {
		t.Errorf("expected source %q, got %q", source, got)
	}
	if got := env.Get("KOOL_WORKSPACE_PATH"); got != workspace {
		t.Errorf("expected workspace path %q, got %q", workspace, got)
	}
	expectedProject := "myapp-workspace-" + composeProjectName(workspaceIdentity(workspace))
	if got := env.Get("COMPOSE_PROJECT_NAME"); got != expectedProject {
		t.Errorf("expected worktree Compose project %q, got %q", expectedProject, got)
	}
}

func TestInitGitWorktreeIgnoresMainWorktree(t *testing.T) {
	root := t.TempDir()
	originalGitWorktreeOutput := gitWorktreeOutput
	gitWorktreeOutput = fakeGitWorktreeOutput(map[string]string{
		"rev-parse --show-toplevel":  root,
		"rev-parse --git-dir":        ".git",
		"rev-parse --git-common-dir": ".git",
	})
	defer func() { gitWorktreeOutput = originalGitWorktreeOutput }()

	env := NewFakeEnvStorage()
	initGitWorktree(env, root)

	if env.IsTrue("KOOL_WORKSPACE") {
		t.Error("expected main Git worktree to remain the source project")
	}
	if got := env.Get("KOOL_WORKSPACE_PROVIDER"); got != "" {
		t.Errorf("expected no workspace provider, got %q", got)
	}
}

func TestInitGitWorktreePreservesRiftContext(t *testing.T) {
	env := NewFakeEnvStorage()
	env.Set("KOOL_WORKSPACE_PROVIDER", "rift")

	originalGitWorktreeOutput := gitWorktreeOutput
	gitWorktreeOutput = func(string, ...string) ([]byte, error) {
		return nil, fmt.Errorf("Git should not be called")
	}
	defer func() { gitWorktreeOutput = originalGitWorktreeOutput }()

	initGitWorktree(env, t.TempDir())

	if got := env.Get("KOOL_WORKSPACE_PROVIDER"); got != "rift" {
		t.Errorf("expected Rift context to be preserved, got %q", got)
	}
}

func TestDuplicateGitWorktreeBasenamesUseUniqueHosts(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "one", "task")
	second := filepath.Join(root, "two", "task")
	for _, workspace := range []string{first, second} {
		if err := os.MkdirAll(workspace, 0755); err != nil {
			t.Fatal(err)
		}
	}
	output := []byte("worktree " + first + "\x00HEAD abc\x00\x00worktree " + second + "\x00HEAD def\x00")
	if !hasDuplicateWorktreeBasename(output, first) || !hasDuplicateWorktreeBasename(output, second) {
		t.Fatal("expected duplicate worktree basenames to require unique hosts")
	}
	if hasDuplicateWorktreeBasename(output, filepath.Join(root, "three", "unique")) {
		t.Fatal("did not expect a unique worktree basename to require a hash")
	}
}

func fakeGitWorktreeOutput(outputs map[string]string) func(string, ...string) ([]byte, error) {
	return func(_ string, args ...string) ([]byte, error) {
		key := strings.Join(args, " ")
		output, found := outputs[key]
		if !found {
			return nil, fmt.Errorf("unexpected Git command: %s", key)
		}
		return []byte(output), nil
	}
}
