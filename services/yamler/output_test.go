package yamler

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestWriteYAML(t *testing.T) {
	o := &DefaultOutputWritter{}
	f := filepath.Join(t.TempDir(), "test.yml")

	y := new(yaml.Node)
	yml := "foo: bar"

	_ = yaml.Unmarshal([]byte(yml), y)

	if err := o.WriteYAML(f, y); err != nil {
		t.Errorf("unexpected error writing YAML: %v", err)
	}

	fh, _ := os.Open(f)
	bs, _ := io.ReadAll(fh)
	fh.Close()
	got := strings.Trim(string(bs), " \t\n")

	if got != yml {
		t.Errorf("bad YML; expected '%s' but got '%s'", yml, got)
	}
}

func TestWriteYAMLIndentation(t *testing.T) {
	o := &DefaultOutputWritter{}
	f := filepath.Join(t.TempDir(), "test_indent.yml")

	y := new(yaml.Node)
	yml := "foo:\n    bar: xxx"

	_ = yaml.Unmarshal([]byte(yml), y)

	if err := o.WriteYAML(f, y); err != nil {
		t.Errorf("unexpected error writing YAML: %v", err)
	}

	fh, _ := os.Open(f)
	bs, _ := io.ReadAll(fh)
	fh.Close()
	got := strings.Trim(string(bs), " \t\n")

	expect := "foo:\n  bar: xxx"

	if got != expect {
		t.Errorf("bad YML; expected '%s' but got '%s'", expect, got)
	}
}
