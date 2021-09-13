//go:build !windows
// +build !windows

package environment

import (
	"fmt"
	"os"
	"testing"
)

func TestInitAsuser(t *testing.T) {
	f := NewFakeEnvStorage()
	initAsuser(f)

	if f.Envs["KOOL_ASUSER"] != fmt.Sprintf("%d", os.Getuid()) {
		t.Error("failed setting current user to KOOL_ASUSER")
	}
}

func TestAlreadyExistingKoolUserInitAsuser(t *testing.T) {
	f := NewFakeEnvStorage()
	f.Envs["KOOL_ASUSER"] = "testing_user"

	initAsuser(f)

	if f.Envs["KOOL_ASUSER"] != "testing_user" {
		t.Error("should not set new user if it is already set")
	}
}
