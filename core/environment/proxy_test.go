package environment

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitProxySourceHost(t *testing.T) {
	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, "kool.yml"), []byte("proxy:\n  domain: exlink.localhost\n  routes: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	env := NewFakeEnvStorage()

	initProxy(env, workDir)

	if got := env.Get("KOOL_PROXY_DOMAIN"); got != "exlink.localhost" {
		t.Errorf("expected proxy domain, got %q", got)
	}
	if got := env.Get("KOOL_PROXY_HOST"); got != "exlink.localhost" {
		t.Errorf("expected source proxy host, got %q", got)
	}
}

func TestInitProxyWorkspaceHost(t *testing.T) {
	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, "kool.yml"), []byte("proxy:\n  domain: exlink.localhost\n  routes: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	env := NewFakeEnvStorage()
	env.Set("KOOL_WORKSPACE", "true")
	env.Set("KOOL_WORKSPACE_NAME", "vite-smoke")

	initProxy(env, workDir)

	if got := env.Get("KOOL_PROXY_HOST"); got != "vite-smoke.workspace.exlink.localhost" {
		t.Errorf("expected workspace proxy host, got %q", got)
	}
}

func TestInitProxyGitWorktreeHost(t *testing.T) {
	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, "kool.yml"), []byte("proxy:\n  domain: exlink.localhost\n  routes: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	env := NewFakeEnvStorage()
	env.Set("KOOL_WORKSPACE", "true")
	env.Set("KOOL_WORKSPACE_NAME", "vite-smoke")
	env.Set("KOOL_WORKSPACE_PROVIDER", "worktree")

	initProxy(env, workDir)

	if got := env.Get("KOOL_PROXY_HOST"); got != "vite-smoke.workspace.exlink.localhost" {
		t.Errorf("expected Git worktree proxy host, got %q", got)
	}
}
