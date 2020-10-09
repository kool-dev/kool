package environment

import (
	"sort"
	"testing"
)

func TestFakeEnvStorage(t *testing.T) {
	f := NewFakeEnvStorage()

	f.Set("testing_key", "testing_value")

	if f.Envs["testing_key"] != "testing_value" {
		t.Error("failed to set a new environment variable on FakeEnvStorage")
	}

	got := f.Get("testing_key")

	if f.Envs["testing_key"] != got {
		t.Errorf("expecting value 'testing_value' on FakeEnvStorage Get, got '%s'", got)
	}

	_ = f.Load("")

	if !f.CalledLoad {
		t.Error("failed to call Load on FakeEnvStorage")
	}
}

func TestAllFakeEnvStorage(t *testing.T) {
	f := NewFakeEnvStorage()

	f.Envs["VAR_1"] = "1"
	f.Envs["VAR_2"] = "2"

	envs := f.All()

	sort.Strings(envs)

	if len(envs) != 2 || envs[0] != "VAR_1=1" || envs[1] != "VAR_2=2" {
		t.Error("failed to return all environment variables")
	}
}

func TestIsTrueUnsetVariableFakeEnvStorage(t *testing.T) {
	f := NewFakeEnvStorage()

	if f.IsTrue("undefined-env-variable") {
		t.Error("Undefined environment variable thought as true.")
	}
}

func TestIsTrueNumeric1FakeEnvStorage(t *testing.T) {
	f := NewFakeEnvStorage()
	f.Envs["env-numeric"] = "1"

	if !f.IsTrue("env-numeric") {
		t.Error("Environment variable with value 1 should be true.")
	}
}

func TestIsTrueStringTrueFakeEnvStorage(t *testing.T) {
	f := NewFakeEnvStorage()
	f.Envs["env-string"] = "true"

	if !f.IsTrue("env-string") {
		t.Error("Environment variable with value 'true' should be true.")
	}
}

func TestIsTrueNonBooleanStringFakeEnvStorage(t *testing.T) {
	f := NewFakeEnvStorage()
	f.Envs["env-other"] = "something"

	if f.IsTrue("env-other") {
		t.Error("Environment variable non-boolean value should not be true.")
	}
}
