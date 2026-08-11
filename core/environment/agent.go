package environment

import "strings"

var agentEnvVars = []string{
	"AI_AGENT",
	"CURSOR_AGENT",
	"GEMINI_CLI",
	"CODEX_SANDBOX",
	"CODEX_CI",
	"CODEX_THREAD_ID",
	"AUGMENT_AGENT",
	"OPENCODE_CLIENT",
	"OPENCODE",
	"AMP_CURRENT_THREAD_ID",
	"CLAUDECODE",
	"CLAUDE_CODE",
	"CLAUDE_CODE_IS_COWORK",
	"REPL_ID",
	"COPILOT_MODEL",
	"COPILOT_ALLOW_ALL",
	"COPILOT_CLI",
	"ANTIGRAVITY_AGENT",
	"PI_CODING_AGENT",
	"KIRO_AGENT_PATH",
}

// AgentEnvVars returns agent environment variables that are safe to forward.
func AgentEnvVars(env EnvStorage) []string {
	values := make(map[string]string)
	for _, entry := range env.All() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			values[parts[0]] = parts[1]
		}
	}

	var forwarded []string
	for _, name := range agentEnvVars {
		if value, exists := values[name]; exists {
			forwarded = append(forwarded, name+"="+value)
		}
	}

	return forwarded
}
