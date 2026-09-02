package automate

import (
	"testing"

	"github.com/spf13/afero"
	"kool-dev/kool/core/shell"
)

func TestReplaceCopiedFiles(t *testing.T) {
	fs := afero.NewMemMapFs()
	e := NewExecutor(&shell.FakeShell{}, func(string) ([]byte, error) {
		return []byte("image: gradle:8-jdkJAVA_VERSION\nport: ${KOOL_APP_PORT:-8080}\n"), nil
	})
	e.local = fs
	e.promptState["JAVA_VERSION"] = "21"
	t.Setenv("JAVA_VERSION", "21")

	if err := e.Do([]*ActionSet{{Actions: []*Action{
		{Src: "docker-compose.yml"},
		{Replace: "JAVA_VERSION $JAVA_VERSION"},
	}}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := afero.ReadFile(fs, "docker-compose.yml")
	if err != nil {
		t.Fatalf("failed reading copied file: %v", err)
	}
	if got, want := string(data), "image: gradle:8-jdk21\nport: ${KOOL_APP_PORT:-8080}\n"; got != want {
		t.Errorf("unexpected file content: got %q, want %q", got, want)
	}
}
