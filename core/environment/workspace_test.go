package environment

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkspaceIdentityDistinguishesDuplicateBasenames(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "one", "task")
	second := filepath.Join(root, "two", "task")
	for _, workspace := range []string{first, second} {
		if err := os.MkdirAll(workspace, 0755); err != nil {
			t.Fatal(err)
		}
	}
	if workspaceIdentity(first) == workspaceIdentity(second) {
		t.Fatal("expected duplicate workspace basenames to have distinct identities")
	}
}

func TestWorkspaceIdentityPreservesHashForLongBasenames(t *testing.T) {
	root := t.TempDir()
	name := strings.Repeat("a", 70)
	first := filepath.Join(root, "one", name)
	second := filepath.Join(root, "two", name)
	for _, workspace := range []string{first, second} {
		if err := os.MkdirAll(workspace, 0755); err != nil {
			t.Fatal(err)
		}
	}
	firstIdentity := workspaceIdentity(first)
	secondIdentity := workspaceIdentity(second)
	if len(firstIdentity) > 63 || len(secondIdentity) > 63 {
		t.Fatalf("expected identities to fit a DNS label: %q %q", firstIdentity, secondIdentity)
	}
	if firstIdentity == secondIdentity {
		t.Fatal("expected long duplicate basenames to preserve distinct hash suffixes")
	}
}

func TestWorkspaceHostDoesNotExposeIdentityHash(t *testing.T) {
	workspace := filepath.Join(t.TempDir(), "quiet-yarrow")
	if err := os.MkdirAll(workspace, 0755); err != nil {
		t.Fatal(err)
	}
	env := NewFakeEnvStorage()
	initWorkspaceContext(env, workspace, filepath.Dir(workspace), workspace, "worktree", true)

	if got := env.Get("KOOL_WORKSPACE_NAME"); got != "quiet-yarrow" {
		t.Fatalf("expected clean workspace host name, got %q", got)
	}
	if project := env.Get("KOOL_WORKSPACE_PROJECT"); !strings.Contains(project, workspaceIdentity(workspace)) {
		t.Fatalf("expected internal project identity to remain unique, got %q", project)
	}
}

func TestIsolateWorkspacePreservesParentOverride(t *testing.T) {
	parent, err := os.CreateTemp("", "kool-workspace-parent-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err = parent.Close(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Remove(parent.Name())
		CleanupWorkspace()
	})
	workspaceOverrideFile = parent.Name()

	restore := IsolateWorkspace()
	child, err := os.CreateTemp("", "kool-workspace-child-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err = child.Close(); err != nil {
		t.Fatal(err)
	}
	workspaceOverrideFile = child.Name()
	restore()

	if _, err = os.Stat(parent.Name()); err != nil {
		t.Fatalf("expected parent override to remain: %v", err)
	}
	if _, err = os.Stat(child.Name()); !os.IsNotExist(err) {
		t.Fatalf("expected child override to be removed, got %v", err)
	}
	if workspaceOverrideFile != parent.Name() {
		t.Fatalf("expected parent override to be restored, got %q", workspaceOverrideFile)
	}
}
