package environment

import (
	"os"
	"path/filepath"
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
