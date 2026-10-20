package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestProxyAliasesCommand(t *testing.T) {
	config := `proxy:
  domain: app.localhost
  network: shared
  routes:
    node:
      ports: ["3001:3001"]
    app:
      ports: ["80:8080", "8081:8081"]
`
	for _, workspace := range []bool{false, true} {
		name := "source"
		if workspace {
			name = "workspace"
		}
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "kool.yml"), []byte(config), 0644); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PWD", dir)
			t.Setenv("KOOL_NAME", "fallback")
			t.Setenv("COMPOSE_PROJECT_NAME", "compose-fallback")
			t.Setenv("KOOL_WORKSPACE_SOURCE_PROJECT", "example")
			t.Setenv("KOOL_WORKSPACE_PROJECT", "example-workspace-task-a")
			t.Setenv("KOOL_PROXY_HOST", "app.localhost")
			t.Setenv("KOOL_WORKSPACE", "false")
			project, alias := "example", "kool-2c16ea1198025827"
			if workspace {
				t.Setenv("KOOL_WORKSPACE", "true")
				project, alias = "example-workspace-task-a", "kool-1c11a4b25a7f4f2e"
			}
			// Listing configured names must not require Docker or Caddy.
			t.Setenv("PATH", t.TempDir())
			root := &cobra.Command{Use: "kool"}
			AddKoolProxy(root)
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&out)
			root.SetArgs([]string{"proxy", "aliases"})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			output := out.String()
			for _, expected := range []string{"PROJECT", "SERVICE", "NETWORK", "ALIAS", project, "shared", "app", "node", alias} {
				if !strings.Contains(output, expected) {
					t.Errorf("expected %q in output:\n%s", expected, output)
				}
			}
			if strings.Count(output, alias) != 1 {
				t.Fatalf("multiple ports must not duplicate a service alias:\n%s", output)
			}
			if strings.Index(output, "app") > strings.Index(output, "node") {
				t.Fatalf("expected services sorted by name:\n%s", output)
			}
		})
	}
}

func TestProxyAliasesCommandEmptyAndInvalidConfig(t *testing.T) {
	for _, test := range []struct {
		name, config, want string
		wantError          bool
	}{
		{"missing config", "", "No proxy aliases configured.", false},
		{"no proxy", "scripts: {}\n", "No proxy aliases configured.", false},
		{"empty routes", "proxy:\n  domain: app.localhost\n", "No proxy aliases configured.", false},
		{"invalid config", "proxy:\n  domain: app.localhost\n  routes:\n    app:\n      ports: [invalid]\n", "listen:target", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			if test.config != "" {
				if err := os.WriteFile(filepath.Join(dir, "kool.yml"), []byte(test.config), 0644); err != nil {
					t.Fatal(err)
				}
			}
			t.Setenv("PWD", dir)
			t.Setenv("PATH", t.TempDir())
			root := &cobra.Command{Use: "kool"}
			AddKoolProxy(root)
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&out)
			root.SetArgs([]string{"proxy", "aliases"})
			err := root.Execute()
			if test.wantError {
				if err == nil || !strings.Contains(err.Error(), test.want) {
					t.Fatalf("expected %q error, got %v", test.want, err)
				}
			} else if err != nil || !strings.Contains(out.String(), test.want) {
				t.Fatalf("expected %q, got %q, error: %v", test.want, out.String(), err)
			}
		})
	}
}

func TestProxyAliasesCommandRejectsArguments(t *testing.T) {
	root := &cobra.Command{Use: "kool"}
	AddKoolProxy(root)
	root.SetArgs([]string{"proxy", "aliases", "app"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected aliases to reject positional arguments")
	}
}
