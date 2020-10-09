package environment

import (
	"os"
	"testing"
)

func TestEnvStorage(t *testing.T) {
	e := NewEnvStorage()

	if _, ok := e.(*DefaultEnvStorage); !ok {
		t.Error("unexpected EnvStorage on NewEnvStorage")
	}

	e.Set("VAR_TESTING_ENV_STORAGE", "1")

	if value, present := os.LookupEnv("VAR_TESTING_ENV_STORAGE"); !present || value != "1" {
		t.Error("failed to set environment variable on EnvStorage")
	}

	os.Setenv("VAR_TESTING_ENV_STORAGE_2", "2")

	if value := e.Get("VAR_TESTING_ENV_STORAGE_2"); value != "2" {
		t.Error("failed to get environment variable on EnvStorage")
	}

	err := e.Load(".env.testing")

	if err != nil {
		t.Errorf("unexpected error to load .env.testing: %v", err)
	}

	if value, present := os.LookupEnv("VAR_TESTING_FILE"); !present || value != "1" {
		t.Error("failed to load environment file on EnvStorage")
	}
}

func TestAllEnvStorage(t *testing.T) {
	e := NewEnvStorage()

	os.Setenv("TESTING_ALL_VAR_1", "1")
	os.Setenv("TESTING_ALL_VAR_1", "1")

	envs := e.All()

	if len(envs) != len(os.Environ()) {
		t.Error("failed to return all environment variables")
	}
}

func TestIsTrueUnsetVariableEnvStorage(t *testing.T) {
	e := NewEnvStorage()

	if e.IsTrue("undefined-env-variable") {
		t.Error("Undefined environment variable thought as true.")
	}
}

func TestIsTrueNumeric1EnvStorage(t *testing.T) {
	e := NewEnvStorage()

	os.Setenv("env-numeric", "1")

	if !e.IsTrue("env-numeric") {
		t.Error("Environment variable with value 1 should be true.")
	}
}

func TestIsTrueStringTrueEnvStorage(t *testing.T) {
	e := NewEnvStorage()

	os.Setenv("env-string", "true")

	if !e.IsTrue("env-string") {
		t.Error("Environment variable with value 'true' should be true.")
	}
}

func TestIsTrueNonBooleanStringEnvStorage(t *testing.T) {
	e := NewEnvStorage()

	os.Setenv("env-other", "something")

	if e.IsTrue("env-other") {
		t.Error("Environment variable non-boolean value should not be true.")
	}
}
