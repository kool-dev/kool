package environment

import (
	"testing"
)

func TestInitAsuser(t *testing.T) {
	f := NewFakeEnvStorage()
	initAsuser(f)

	if f.Envs["KOOL_ASUSER"] != uid() {
		t.Error("failed setting current user to KOOL_ASUSER")
	}
}

func TestAlreadyExistingKoolUserInitAsuser(t *testing.T) {
	f := NewFakeEnvStorage()
	f.Envs["KOOL_ASUSER"] = "testing_user"

	initAsuser(f)

	if f.Envs["KOOL_ASUSER"] != uid() {
		t.Error("failed setting current user to KOOL_ASUSER")
	}
}
