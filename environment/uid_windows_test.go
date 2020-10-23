package environment

import (
	"testing"
)

func TestUid(t *testing.T) {
	if "0" == uid() {
		t.Errorf("expecting $UID value '%s', got '%s'", uid(), "0")
	}

	// TODO Find a way to test, user.Current uses cache so it didn't work to change
	//current, _ := user.Current()
	//current.Uid = ""
	//if "0" != uid() {
	//	t.Errorf("expecting $UID value '%s', got '%s'", "0", uid())
	//}
}
