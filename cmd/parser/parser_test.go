package parser

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestDefaultParser(t *testing.T) {
	var p Parser = NewParser()

	if _, assert := p.(*DefaultParser); !assert {
		t.Errorf("NewParser() did not return a *DefaultParser")
	}
}

func TestParserAddLooupPath(t *testing.T) {
	var (
		p      Parser = NewParser()
		err    error
		tmpDir = t.TempDir()
	)

	err = p.AddLookupPath(tmpDir)

	if err == nil || ErrKoolYmlNotFound.Error() != err.Error() {
		t.Errorf("expected ErrKoolYmlNotFound; got %s", err)
		return
	}

	if err = ioutil.WriteFile(path.Join(tmpDir, "kool.yml"), []byte(KoolYmlOK), os.ModePerm); err != nil {
		t.Fatalf("failed creating temp file; error: %s", err)
	}

	if err = p.AddLookupPath(tmpDir); err != nil {
		t.Errorf("unexpected error; error: %s", err)
		return
	}
}
