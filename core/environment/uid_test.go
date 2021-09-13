//go:build !windows
// +build !windows

package environment

import (
	"fmt"
	"os"
	"testing"
)

func TestInitUid(t *testing.T) {
	f := NewFakeEnvStorage()
	initUid(f)

	if f.Envs["UID"] != fmt.Sprintf("%d", os.Getuid()) {
		t.Error("failed setting current uid to UID")
	}
}

func TestAlreadyExistingKoolUserInitUid(t *testing.T) {
	f := NewFakeEnvStorage()
	f.Envs["UID"] = "1000"

	initUid(f)

	if f.Envs["UID"] != "1000" {
		t.Error("should not set new uid if it is already set")
	}
}
