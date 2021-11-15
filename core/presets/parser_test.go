package presets

import (
	"embed"
	"testing"

	"github.com/leaanthony/debme"
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
