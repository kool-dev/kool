package commands

import (
	"kool-dev/kool/core/shell"
	"testing"
)

func TestKoolServiceProxies(t *testing.T) {
	k := &DefaultKoolService{
		&shell.FakeShell{},
	}

	if _, ok := k.Shell().(*shell.FakeShell); !ok {
		t.Error("unexpected Shell return")
	}
}
