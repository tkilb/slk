// internal/config/fork_grammar.go
//
// Fork feature grammar-check: the [grammar] config section. Every field
// is optional; zero values fall back to the defaults in internal/grammar.
package config

// Grammar configures Ctrl+G, the compose-box proofreader that runs the
// `claude` CLI.
type Grammar struct {
	// Model is passed to `claude --model` (default "haiku").
	Model string `toml:"model"`
	// TimeoutSeconds bounds one request (default 30).
	TimeoutSeconds int `toml:"timeout_seconds"`
	// Prompt replaces the built-in proofreader system prompt.
	Prompt string `toml:"prompt"`
	// Command is the claude executable (default "claude", from PATH).
	Command string `toml:"command"`
}
