# Flag Diff: Open Delivery vs Wukong

- generated_at: `2026-07-07T06:32:25Z`
- scope: visible leaf commands only; parses local `Flags:` from `dws <path> --help`; excludes global flags, hidden flags, and Cobra automatic `--help`.
- open_leaf_commands: `631`
- wukong_leaf_commands: `781`
- common_leaf_commands: `626`
- open_help_errors: `0`
- wukong_help_errors: `0`

## Summary

- common commands where open has extra visible flags: `4`
- common commands where wukong has extra visible flags: `5`
- common commands with changed visible flag value type: `0`
- open-only leaf commands with local flags: `5`

## Open Extra Visible Flags On Common Commands

- `chat message reply`: `--ai-tag`
- `chat message send`: `--ai-tag`
- `pat chmod`: `--domain`, `--domains`, `--product`, `--products`, `--recommend`

## Wukong Extra Visible Flags On Common Commands

- `chat message list`: `--forward`
- `chat message list-topic-replies`: `--forward`
- `chat message send`: `--dentry-id`, `--file-name`, `--file-size`, `--file-type`, `--space-id`
- `skill install`: `--force`, `--skill-id`

## Open-Only Leaf Command Flags

- `chat media upload`: `--file`, `--type`
- `doctor`: `--json`, `--perf`, `--timeout`
- `pat browser-policy`: `--agentCode`, `--enabled`
- `skill get`: `--skill-id`
- `upgrade`: `--all`, `--check`, `--force`, `--list`, `--rollback`, `--skip-skills`, `--version`

## Changed Flag Types

- none

## Compatibility Notes

- `--ai-tag` is an open-source customization on `chat message send` and `chat message reply`. These are the only two visible chat leaf commands that pass `clawType` into `send_personal_message`; no additional `--ai-tag` leaf was found missing.
- `chat message list --forward`, `chat message list-topic-replies --forward`, and `chat message send --dentry-id/--space-id/--file-name/--file-type/--file-size` are intentionally hidden in open help, but remain registered and executable for old command compatibility.
- `skill search` is now aligned with Wukong help: visible `--source`; deprecated `--scopes` is kept as hidden compatibility and normalized to the same `source=` request parameter.
- `skill install --skill-id/--force` belongs to the Wukong App install flow. Open-source main/v1.0.47 already used positional `dws skill install <skillId> <target>` for Agent directory installs, so this is an edition difference rather than an open main regression.
