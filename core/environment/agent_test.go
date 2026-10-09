package environment

import (
	"reflect"
	"testing"
)

func TestAgentEnvVars(t *testing.T) {
	env := NewFakeEnvStorage()
	env.Envs = map[string]string{
		"AI_AGENT":             "custom",
		"OPENCODE":             "1",
		"CLAUDECODE":           "",
		"COPILOT_GITHUB_TOKEN": "secret",
		"UNRELATED":            "value",
	}

	expected := []string{"AI_AGENT=custom", "OPENCODE=1", "CLAUDECODE="}
	if actual := AgentEnvVars(env); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
