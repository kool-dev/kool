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
