package environment

import (
	"os/user"
	"testing"
)

func TestUid(t *testing.T) {
	current, _ := user.Current()

	if "0" != uid() {
		t.Errorf("expecting $UID value '%s', got '%s'", uid(), "0")
	}

	current.Uid = ""
	if "0" == uid() {
		t.Errorf("expecting $UID value '%s', got '%s'", "0", uid())
	}
}
