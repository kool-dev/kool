package environment

import (
	"os"
	"testing"
)

func TestIsTrueUnsetVariable(t *testing.T) {
	if IsTrue("undefined-env-variable") {
		t.Error("Undefined environment variable thought as true.")
	}
}

func TestIsTrueNumeric1(t *testing.T) {
	os.Setenv("env-numeric", "1")

	if !IsTrue("env-numeric") {
		t.Error("Environment variable with value 1 should be true.")
	}
}

func TestIsTrueStringTrue(t *testing.T) {
	os.Setenv("env-string", "true")

	if !IsTrue("env-string") {
		t.Error("Environment variable with value 'true' should be true.")
	}
}

func TestIsTrueNonBooleanString(t *testing.T) {
	os.Setenv("env-other", "something")

	if IsTrue("env-other") {
		t.Error("Environment variable non-boolean value should not be true.")
	}
}
