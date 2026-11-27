package commands

import (
	"io"
	"kool-dev/kool/core/environment"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// newTestCloudCommand wires a cloud command with a no-op subcommand under a
// fake root, as its PersistentPreRunE calls the root's one
func newTestCloudCommand() (root *cobra.Command, cloud *Cloud) {
	cloud = NewCloud()
	cloud.env = environment.NewFakeEnvStorage()

	cloudCmd := NewCloudCommand(cloud)
	cloudCmd.AddCommand(&cobra.Command{Use: "deploy", RunE: func(*cobra.Command, []string) error { return nil }})

	root = &cobra.Command{Use: "kool", PersistentPreRunE: func(*cobra.Command, []string) error { return nil }}
	root.AddCommand(cloudCmd)
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)

	return
}

func TestCloudCommandHelpPrefersTokenFromEnv(t *testing.T) {
	cmd := NewCloudCommand(NewCloud())

	if !strings.Contains(cmd.Long, "KOOL_API_TOKEN") || !strings.Contains(cmd.Long, "process list") {
		t.Errorf("expected cloud help to say the token is read from KOOL_API_TOKEN, preferred over --token; got %q", cmd.Long)
	}

	token := cmd.PersistentFlags().Lookup("token")
	if token == nil {
		t.Fatal("expected the --token flag to still exist")
	}

	if !strings.Contains(token.Usage, "KOOL_API_TOKEN") || !strings.Contains(token.Usage, "process list") {
		t.Errorf("expected --token help to point to KOOL_API_TOKEN; got %q", token.Usage)
	}
}

func TestCloudCommandToken(t *testing.T) {
	root, cloud := newTestCloudCommand()
	root.SetArgs([]string{"cloud", "deploy", "--domain", "app.example.com"})

	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "KOOL_API_TOKEN environment variable") {
		t.Errorf("expected missing token error pointing to KOOL_API_TOKEN; got %v", err)
	}

	// the token is read from the environment
	cloud.env.Set("KOOL_API_TOKEN", "env-token")
	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error with KOOL_API_TOKEN set: %v", err)
	}

	// --token keeps working, and wins over the environment
	root, cloud = newTestCloudCommand()
	cloud.env.Set("KOOL_API_TOKEN", "env-token")
	root.SetArgs([]string{"cloud", "deploy", "--domain", "app.example.com", "--token", "flag-token"})

	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error with --token: %v", err)
	}

	if token := cloud.env.Get("KOOL_API_TOKEN"); token != "flag-token" {
		t.Errorf("expected --token to be shared as KOOL_API_TOKEN; got %q", token)
	}
}
