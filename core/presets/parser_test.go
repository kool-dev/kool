package presets

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"kool-dev/kool/core/shell"

	"github.com/leaanthony/debme"
	"gopkg.in/yaml.v2"
)

//go:embed fixtures/*
var fixtures embed.FS

func TestPresetParser(t *testing.T) {
	root, _ := debme.FS(fixtures, "fixtures")

	SetSource(root)

	m := make(map[string]bool)
	p := NewParser()

	for _, t := range p.GetTags() {
		m[t] = true
	}

	if !m["foo"] {
		t.Errorf("missing tag 'foo'; %+v", p.GetTags())
	}

	if !p.Exists("foo") {
		t.Error("preset 'foo' should exist")
	} else if p.Exists("bar") {
		t.Error("preset 'bar' should not exist")
	}

	if len(p.GetPresets("foo")) != 1 {
		t.Error("should have found 1 preset with tag foo")
	}

	if len(p.GetPresets("bar")) != 0 {
		t.Error("should NOT have found any preset with tag bar")
	}
}

func TestBunPresetInstall(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	previousSource := source
	t.Cleanup(func() {
		SetSource(previousSource)
		if err := os.Chdir(cwd); err != nil {
			t.Error(err)
		}
	})
	SetSource(os.DirFS(filepath.Join(cwd, "../..")).(SourceFS))
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	p := NewParser()
	p.PrepareExecutor(&shell.FakeShell{})
	if err := p.Install("bunjs"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("kool.yml")
	if err != nil {
		t.Fatal(err)
	}
	var config struct {
		Scripts map[string]interface{} `yaml:"scripts"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	setup, ok := config.Scripts["setup"].([]interface{})
	if !ok || len(setup) != 2 || setup[1] != "kool start" {
		t.Fatalf("expected install/start setup script, got %#v", config.Scripts["setup"])
	}
}
