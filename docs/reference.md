# CLI Reference

This page documents stable cross-command behavior. Command paths and flags are owned by the executable command tree; inspect them with `--help` instead of relying on a copied command catalog.

## Command discovery

```bash
dws --help
dws <service> --help
dws <service> <command> --help
```

The packaged Agent Skills are the source of task-oriented guidance for agents. `dws --help` is the source of truth for the command surface installed on the current machine.

## Configuration

Use the built-in registry to list supported environment variables, defaults, examples, and their current values. Sensitive values are masked.

```bash
dws config list
dws config list --json
```

Common registered overrides include `DWS_CONFIG_DIR`, `DWS_CLIENT_ID`, `DWS_CLIENT_SECRET`, `DWS_LANG`, and `DWS_TRUSTED_DOMAINS`. Use the registry instead of maintaining a duplicate environment-variable table.

On macOS only, `DWS_DISABLE_KEYCHAIN=1` is an emergency fallback for sandboxed runtimes that cannot access Keychain APIs. It stores the encryption key beside the encrypted token and therefore weakens protection at rest. Agent hosts should use the declaration described in [Agent identification](./integrations/agent-identification.md) instead of generic environment detection.

## Output formats

All regular commands share `--format` (`-f`):

| Format | Contract |
|---|---|
| `json` | Structured JSON for agents and scripts. |
| `table` | Human-readable tabular output. |
| `raw` | Unmodified upstream payload where supported. |
| `pretty` | Human-oriented annotated output. |
| `ndjson` | One JSON object per line for streaming and pipelines. |
| `csv` | Tabular CSV output when the response has rows. |

Use `--fields` for field projection and `--jq` for JSON filtering:

```bash
dws contact user search --query "Alice" --fields name,userId
dws contact user search --query "Alice" --jq '.result[] | {name, userId}'
```

Write the rendered result to a file with `--output` (`-o`):

```bash
dws contact user search --query "Alice" --format json -o result.json
```

Event streams default to `ndjson`; bounded streams may also use `json` or `pretty`. Consult the event command's help for stream-specific restrictions.

## Safe execution

`--dry-run` previews a request without sending it. Mutating commands that require confirmation accept `--yes` only after the caller has reviewed the target and parameters.

```bash
dws todo task create --title "Review release" --executors USER_ID --dry-run
dws todo task create --title "Review release" --executors USER_ID --yes
```

## Exit codes

| Code | Category | Meaning |
|---|---|---|
| 0 | Success | The command completed successfully. |
| 1 | API | An MCP tool call or upstream API failed. |
| 2 | Auth | Authentication or authorization failed. |
| 3 | Validation | Input, flags, or parameter validation failed. |
| 4 | PAT | PAT authorization interception; stderr contains the machine-readable PAT payload. |
| 5 | Internal | An unexpected internal error occurred. |
| 6 | Discovery | Endpoint discovery or protocol negotiation failed. |

With `--format json`, structured errors include machine-readable fields such as `category`, `reason`, `hint`, and `actions`.

## Schema introspection

Product commands are compiled into the binary. Use their `--help` output for parameters. `dws schema` remains available only for helper-owned schemas such as `dev.*`.

```bash
dws schema "dev app create"
dws schema --cli-path "dev app create"
dws schema "dev app create" --jq '.tool.parameters'
```

## Shell completion

```bash
# Bash
dws completion bash > /etc/bash_completion.d/dws

# Zsh
dws completion zsh > "${fpath[1]}/_dws"

# Fish
dws completion fish > ~/.config/fish/completions/dws.fish
```
