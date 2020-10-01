// +build !windows

package environment

import (
	"fmt"
	"os"
	"testing"
)

func TestInitAsuser(t *testing.T) {
	oldUser := os.Getenv("KOOL_ASUSER")

	defer func(oldUser string) {
		if oldUser != "" {
			os.Setenv("KOOL_ASUSER", oldUser)
		} else {
			os.Unsetenv("KOOL_ASUSER")
		}
	}(oldUser)

	os.Setenv("KOOL_ASUSER", "")

	initAsuser()

	if os.Getenv("KOOL_ASUSER") != fmt.Sprintf("%d", os.Getuid()) {
		t.Error("failed setting current user to KOOL_ASUSER")
	}
}

func TestAlreadyExistingKoolUserInitAsuser(t *testing.T) {
	oldUser := os.Getenv("KOOL_ASUSER")

	defer func(oldUser string) {
		if oldUser != "" {
			os.Setenv("KOOL_ASUSER", oldUser)
		} else {
			os.Unsetenv("KOOL_ASUSER")
		}
	}(oldUser)

	os.Setenv("KOOL_ASUSER", "testing_user")

	initAsuser()

	if os.Getenv("KOOL_ASUSER") != "testing_user" {
		t.Error("should not set new user if it is already set")
	}
}
