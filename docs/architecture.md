# Architecture

`dws` is a Go CLI that exposes DingTalk capabilities to humans and AI agents through a compiled Cobra command tree. Product commands use generated static endpoint metadata; helper-owned surfaces such as `dev` are implemented directly in Go. Optional plugins can add commands without changing the core binary.

## High-Level Flow

1. `cmd` is the CLI entrypoint, invoking `internal/app` to build the root Cobra command tree.
2. `internal/app` wires utility commands, generated product commands, helper-owned command groups, and optional plugins.
3. `internal/syncdata` supplies the generated static endpoint and routing data used to build product commands.
4. `internal/helpers` implements helper-owned and customized product behavior such as `dev`, `chat`, `sheet`, and connector workflows.
5. `internal/executor` and `internal/transport` execute MCP JSON-RPC calls; `internal/output` formats responses.
6. `internal/auth` manages login state, PAT tokens, and agent-code detection.

## Repository Structure

- `cmd`: CLI entrypoint
- `internal/app`: root command wiring, static utility commands, and plugin loading
- `internal/helpers`: product command handlers (dev, chat, calendar, contact, etc.)
- `internal/plugin`: optional plugin discovery and command conversion
- `internal/cli`: catalog types and static endpoint loading
- `internal/executor`: invocation dispatch and result handling
- `internal/transport`: MCP HTTP client and request signing
- `internal/auth`: login, token management, agent-code detection, identity
- `internal/audit`: user operation audit log (JSONL, hash chain, forwarding)
- `internal/errors`: structured error model with categories and hints
- `internal/keychain`: OS keychain integration for credential storage
- `internal/security`: endpoint allowlist and domain trust
- `internal/safety`: runtime safety checks (confirm prompts, dry-run guards)
- `internal/cobracmd`: shared Cobra command builders
- `internal/pat`: PAT (Personal Access Token) authorization flow
- `internal/output`: response formatting (json, table, raw, pretty)
- `internal/logging`: structured logging and argument sanitization
- `internal/tui`: terminal UI helpers
- `internal/recovery`: panic recovery and graceful degradation
- `pkg/configmeta`: environment variable registry and documentation
- `pkg/config`: configuration constants and paths
- `pkg/edition`: edition detection (oss vs enterprise)
- `pkg/mcptypes`: MCP protocol type definitions
- `internal/syncdata`: generated static endpoint and command-routing data synced from the Wukong baseline
- `skills/`: bundled agent skills (mono/ and multi/ layouts)
- `test/`: CLI, integration, contract, unit, and skill E2E tests
- `scripts/`: install scripts, policy checks, and CI helpers
