package parser

import (
	"io/ioutil"
	"kool-dev/kool/cmd/builder"
	"os"
	"path"
	"testing"
)

const KoolYmlOK = `scripts:
  single-line: single line script
  multi-line:
    - line 1
    - line 2
`

func TestParseKoolYaml(t *testing.T) {
	var (
		err     error
		tmpPath string
		parsed  *KoolYaml
		cmds    []*builder.Command
	)

	tmpPath = path.Join(t.TempDir(), "kool.yml")

	err = ioutil.WriteFile(tmpPath, []byte(KoolYmlOK), os.ModePerm)

	if err != nil {
		t.Fatal("failed creating temporary file for test", err)
	}

	parsed, err = ParseKoolYaml(tmpPath)

	if err != nil {
		t.Errorf("failed parsing proper kool.yml file; error: %s", err)
		return
	}

	if len(parsed.Scripts) != 2 {
		t.Errorf("expected to parse 2 scripts; got %d", len(parsed.Scripts))
		return
	}

	if !parsed.HasScript("single-line") || !parsed.HasScript("multi-line") {
		t.Errorf("expected to have single-line and multi-line script")
		return
	}

	if cmds, err = parsed.ParseCommands("single-line"); err != nil {
		t.Errorf("failed to parse proper single-line; error: %s", err)
		return
	}

	if len(cmds) != 1 {
		t.Errorf("expected single-line to parse 1 command; got %d", len(cmds))
		return
	}

	if cmds, err = parsed.ParseCommands("multi-line"); err != nil {
		t.Errorf("failed to parse proper multi-line; error: %s", err)
		return
	}

	if len(cmds) != 2 {
		t.Errorf("expected multi-line to parse 1 command; got %d", len(cmds))
		return
	}
}
