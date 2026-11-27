package api

import (
	"kool-dev/kool/core/environment"
	"strings"
	"testing"
)

func TestDetectClient(t *testing.T) {
	tests := []struct {
		name     string
		envs     map[string]string
		expected string
	}{
		{"no markers", map[string]string{}, "human"},
		{"claude code", map[string]string{"CLAUDECODE": "1"}, "claude-code"},
		{"claude code wins over AI_AGENT and CI", map[string]string{"CLAUDECODE": "1", "AI_AGENT": "other", "CI": "true"}, "claude-code"},
		{"CLAUDECODE other than 1", map[string]string{"CLAUDECODE": "0"}, "human"},
		{"AI_AGENT", map[string]string{"AI_AGENT": "pi"}, "pi"},
		{"AI_AGENT sanitized", map[string]string{"AI_AGENT": " My Agent/v1.2_beta-3!\n"}, "myagentv1.2_beta-3"},
		{"AI_AGENT truncated", map[string]string{"AI_AGENT": strings.Repeat("a", 50)}, strings.Repeat("a", 40)},
		{"AI_AGENT empty after sanitizing", map[string]string{"AI_AGENT": "!!!", "CI": "1"}, "ci"},
		{"AI_AGENT wins over documented markers", map[string]string{"AI_AGENT": "pi", "GEMINI_CLI": "1"}, "pi"},
		{"gemini cli", map[string]string{"GEMINI_CLI": "1"}, "gemini-cli"},
		{"copilot cli", map[string]string{"COPILOT_CLI": "1"}, "copilot-cli"},
		{"cursor", map[string]string{"CURSOR_AGENT": "1"}, "cursor"},
		{"agent wins over CI", map[string]string{"GEMINI_CLI": "1", "GITHUB_ACTIONS": "true"}, "gemini-cli"},
		{"CI", map[string]string{"CI": "true"}, "ci"},
		{"GITHUB_ACTIONS", map[string]string{"GITHUB_ACTIONS": "true"}, "ci"},
		{"CI opted out", map[string]string{"CI": "false"}, "human"},
		{"CI zero", map[string]string{"CI": "0"}, "human"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := environment.NewFakeEnvStorage()
			for k, v := range tt.envs {
				env.Set(k, v)
			}

			if client := detectClient(env); client != tt.expected {
				t.Errorf("expected client %q; got %q", tt.expected, client)
			}
		})
	}
}
