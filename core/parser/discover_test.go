package parser

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFindKoolYamlPrefersYml(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"kool.yml", "kool.yaml"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("scripts: {}\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	file, err := FindKoolYaml(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expected := filepath.Join(dir, "kool.yml"); file != expected {
		t.Errorf("expected %q, got %q", expected, file)
	}
}

func TestFindKoolYamlFallsBackToYaml(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "kool.yaml"), []byte("scripts: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	file, err := FindKoolYaml(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expected := filepath.Join(dir, "kool.yaml"); file != expected {
		t.Errorf("expected %q, got %q", expected, file)
	}
}

func TestFindKoolYamlNotFound(t *testing.T) {
	if _, err := FindKoolYaml(t.TempDir()); !errors.Is(err, ErrKoolYmlNotFound) {
		t.Errorf("expected ErrKoolYmlNotFound, got %v", err)
	}
}

func TestLoadKoolYaml(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "kool.yml"), []byte("workspaces:\n  - app\n"), 0644); err != nil {
		t.Fatal(err)
	}

	parsed, err := LoadKoolYaml(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(parsed.Workspaces) != 1 || parsed.Workspaces[0] != "app" {
		t.Errorf("expected workspaces [app], got %v", parsed.Workspaces)
	}
}

func TestLoadKoolYamlMissing(t *testing.T) {
	if _, err := LoadKoolYaml(t.TempDir()); !errors.Is(err, ErrKoolYmlNotFound) {
		t.Errorf("expected ErrKoolYmlNotFound, got %v", err)
	}
}

// A malformed config must be distinguishable from a missing one, so callers
// can choose to report it rather than silently behaving as if unconfigured.
func TestLoadKoolYamlMalformed(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "kool.yml"), []byte("scripts: [oops\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadKoolYaml(dir)
	if err == nil {
		t.Fatal("expected a decoding error, got none")
	}

	if errors.Is(err, ErrKoolYmlNotFound) {
		t.Error("malformed config reported as not found")
	}
}
