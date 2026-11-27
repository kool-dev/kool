package api

import (
	"kool-dev/kool/core/environment"
	"strings"
)

// maxClientLength limits the size of a client name taken from AI_AGENT
const maxClientLength = 40

// detectClient tells who is running the CLI, so the API can know which
// deploys were run by a person, a CI pipeline or a coding agent.
//
// Besides CLAUDECODE and the generic AI_AGENT, it only relies on markers
// each agent documents setting for the commands it runs:
//   - GEMINI_CLI=1: https://github.com/google-gemini/gemini-cli/blob/main/docs/tools/shell.md
//   - COPILOT_CLI=1: https://github.com/github/copilot-cli/blob/main/changelog.md
//   - CURSOR_AGENT: https://cursor.com/docs/agent/tools/terminal
func detectClient(env environment.EnvStorage) string {
	if env.Get("CLAUDECODE") == "1" {
		return "claude-code"
	}

	if client := sanitizeClient(env.Get("AI_AGENT")); client != "" {
		return client
	}

	if env.Get("GEMINI_CLI") == "1" {
		return "gemini-cli"
	}

	if env.Get("COPILOT_CLI") == "1" {
		return "copilot-cli"
	}

	if env.Get("CURSOR_AGENT") != "" {
		return "cursor"
	}

	if isEnvSet(env.Get("CI")) || isEnvSet(env.Get("GITHUB_ACTIONS")) {
		return "ci"
	}

	return "human"
}

// sanitizeClient lowercases the given value and keeps only [a-z0-9._-],
// up to maxClientLength characters
func sanitizeClient(value string) string {
	var client strings.Builder

	for _, r := range strings.ToLower(value) {
		if client.Len() == maxClientLength {
			break
		}

		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			client.WriteRune(r)
		}
	}

	return client.String()
}

// isEnvSet tells whether a flag-like variable is set, not counting
// explicit opt-outs such as CI=false
func isEnvSet(value string) bool {
	return value != "" && value != "0" && strings.ToLower(value) != "false"
}
