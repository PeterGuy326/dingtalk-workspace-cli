# delivery/remove-discovery-yolo Validation Report

Date: 2026-07-07

Fork delivery branch:

```text
https://github.com/PeterGuy326/dingtalk-workspace-cli/tree/delivery/remove-discovery-yolo
```

Release:

```text
https://github.com/PeterGuy326/dingtalk-workspace-cli/releases/tag/delivery-remove-discovery-yolo-latest
```

Install:

```sh
curl -fsSL https://raw.githubusercontent.com/PeterGuy326/dingtalk-workspace-cli/delivery-remove-discovery-yolo-latest/scripts/install-delivery-remove-discovery-yolo.sh | sh
```

## Build

Release binary checked from the `delivery-remove-discovery-yolo-latest` tag target:

```json
{
  "architecture": "MCP Static Endpoint Mode",
  "commit": "<release tag target>",
  "edition": "open",
  "version": "v1.0.50-delivery-remove-discovery-yolo"
}
```

Release assets uploaded:

- `dws-darwin-arm64.tar.gz`
- `dws-darwin-amd64.tar.gz`
- `dws-linux-arm64.tar.gz`
- `dws-linux-amd64.tar.gz`
- `dws-skills.zip`
- `checksums.txt`

## Checks

| Check | Result |
| --- | --- |
| Release curl install into temp HOME/bin | Pass |
| `dws version` shows static endpoint mode | Pass |
| Root `--help` hides retired `cache` and `conference` | Pass |
| `dws cache refresh` compatibility stub | Pass, hidden; points to `dws upgrade` + `internal/syncdata` endpoint check |
| `dws conference` unavailable guidance | Pass |
| `dws dev connect --help` exposes `--yolo`, `--agent-permission-mode`, `--agent-approval-mode` | Pass |
| `dws dev connect --help` does not expose `--agent-yolo` | Pass |
| `dws dev connect` default dry-run yolo | Pass |
| `--agent-permission-mode ask` disables yolo in dry-run | Pass |
| `--agent-approval-mode ask` disables yolo in dry-run | Pass |
| `--agent-yolo` returns unknown flag | Pass |
| Old chat hidden compatibility flags `--sender`, `--senders`, `--keyword`, `--size` | Pass |
| Timestamp compatibility: seconds and milliseconds | Pass |
| `dws skill setup --mode mono --target codex` in temp HOME | Pass |
| `dws skill setup --mode multi --target codex -s dev -s chat` in temp HOME | Pass |
| Recursive visible `dws <path> --help` scan for open delivery | Pass, 779 commands, 0 help errors |
| `dws skill search --help` aligns with Wukong `--source`; old `--scopes` remains hidden/deprecated | Pass |
| Chat old file/forward flags remain registered but hidden from help | Pass |
| PR #562 chat sync: `chat data-auth cross-org` and `message send --msg-type audio/video` | Pass |
| `StaticServers` and `SupplementServers` endpoint injection merge semantics | Pass |
| Source build from a single repository checkout | Pass, static endpoint data is committed under `internal/syncdata` |
| `doc version revert` version preflight uses backend-safe page size | Pass, `maxResults=50` |
| `chat message send` default output has no `[debug]` leakage | Pass |
| `aitable dashboard create/update` unsafe partial paths | Guarded in CLI: `create --name`, `update --name`, and update with only `dashboardName` fail fast |
| Skill sync from Wukong open-source baseline | Pass, 294 files; only same-path diffs are the two chat overlay files |

## Skill Setup Verification

Mono mode:

```text
HOME=<tmp> dws skill setup --mode mono --target codex --yes
installed: <tmp>/.codex/skills/dws/SKILL.md
file count: 133
verified chat reference includes `chat data-auth cross-org` and `--msg-type image/file/audio/video`
```

Multi mode:

```text
HOME=<tmp> dws skill setup --mode multi --target codex --yes -s dev -s chat
installed:
  <tmp>/.codex/skills/dws-shared/SKILL.md
  <tmp>/.codex/skills/dingtalk-chat/SKILL.md
  <tmp>/.codex/skills/dingtalk-dev/SKILL.md
file count: 21
verified chat reference includes `chat data-auth cross-org` and `--msg-type image/file/audio/video`
```

The installed mono skill starts with valid YAML frontmatter and carries the static-endpoint guidance:

```text
命令以当前 dws 二进制为准。服务发现和动态 schema 已下线，本文档随版本内嵌发布；
执行前用 `dws <cmd> --help` 或 `--dry-run` 验证 flag 与命令是否存在。
```

## Command Tree Snapshots

Visible command tree snapshots are committed next to this report:

- `command-tree-open-delivery.md`
- `command-tree-open-delivery.json`
- `command-tree-wukong-local.md`
- `command-tree-wukong-local.json`
- `command-tree-comparison.md`
- `flag-diff-open-vs-wukong.md`
- `flag-diff-open-vs-wukong.json`

Summary:

| Tree | Visible commands | Help errors |
| --- | ---: | ---: |
| Open delivery | 779 | 0 |
| Wukong local build | 973 | 0 |

Skill/help drift notes:

- No install drift: mono and multi skill setup both install successfully in an isolated HOME.
- No command-help reachability drift: every visible open delivery command path has a working `dws <path> --help`.
- Flag drift check: `--ai-tag` is open-only visible customization on `chat message send/reply`; old chat compatibility flags remain hidden but executable; `skill search` now uses visible `--source` and hidden deprecated `--scopes`.
- Content drift found in mono skill coverage: open help exposes `live` and `pat`, but mono skill does not currently include product reference pages for them. Multi mode has `dingtalk-live`; `pat` is still not covered by product skill docs.
- Wukong skill baseline sync: current `skills` has 294 files. Compared with `dws-wukong/target/open-source-cli/skills`, same-path content diffs are limited to `mono/references/products/chat.md` and `multi/dingtalk-chat/references/chat.md`, both expected open overlays for `chat data-auth cross-org` and audio/video msg-type aliases.
- Wukong-only `dingtalk-aiapp` is intentionally not shipped because `aiapp` is not visible in the current open command tree. Open-only `dws-shared`, `dingtalk-dev`, `dingtalk-profile`, and mono `dev.md` are retained.

Top-level wukong-only visible products:

```text
agoal, aiapp, aidesign, blackboard, conference, credit, docparse, finance, law, yida
```

Top-level open-only utilities:

```text
doctor, upgrade
```

Notes:

- The open delivery help surface intentionally excludes retired/internal products from `--help`; old compatibility paths are hidden and verified separately.
- PR #562 review result: current delivery branch already contained `chat data-auth cross-org` as a superset (`--target-org-id` or `--all`). The missing part was `audio`/`video` aliases for `chat message send`; this delivery now accepts both and normalizes them to backend `msgType=file`.
- The `dws cache` command tree is retained only as a hidden/deprecated compatibility entry. Because discovery cache refresh is gone, endpoint miss guidance is now: upgrade first, then verify `internal/syncdata.StaticServers()` coverage.
- Static endpoint and routing data is generated from the Wukong baseline but committed into this repository under `internal/syncdata`; CI and local source builds no longer require any sibling data repository.
- `dws-wukong origin/develop` at `2a7c70ff` currently does not compile directly against the remove-discovery core because it still references removed `pkg/cmdutil` overlay APIs: `IsEnvelopeSourced` and `MergeHardcodedLeaves`.
- The wukong command tree snapshot was generated from the local wukong workspace at `f48ae79d` with the local compatibility bridge changes present in that workspace.
- This delivery branch does not modify `dws-wukong`; wukong is used only as the source baseline for comparison.

## Remaining Failure Triage

Source: `/Users/huyz/Documents/当前仍失败指令清单.md` and the wukong01 `cli_to_mcp` full run on this branch.

Fixed or guarded in this delivery:

- `doc version revert`: previous fake-success guard is kept; the version preflight now calls `list_doc_versions` with `maxResults=50` instead of the backend-invalid `100`.
- `aitable dashboard update --name`: guarded in CLI because partial rename can make a real dashboard unreachable. `dashboard create --name` is also guarded because the backend currently ignores the shortcut and can return an unusable id.
- `chat message send`: non-debug runs no longer write `[debug]` normalization lines.

Needs backend/tool registration decision:

- 11 commands are registered in CLI and have static endpoints, but live `tools/list` does not expose the matching backend tools: `sheet chart list/create/update/delete`, `sheet batch-update`, `sheet range batch-clear`, `todo task add/list/remove-attachment`, `calendar event list-mine`, `minutes record start`.
- Recommendation: either register backend tools, or explicitly mark these like `contact label` as unavailable/hidden. `minutes record start` has a likely tool-name mismatch: CLI calls `execute_listening_note_command`, while live tools show a Chinese-named equivalent.

Needs backend/permission confirmation:

- `mail template update`: CLI source matches the wukong implementation. Wukong skill says `template update` is only valid for draft templates; if a template created with `--is-draft` still returns `Invalid parameter`, treat it as backend or gateway behavior, not remove-discovery routing.
- Chat permission/business failures (`chmod`, `message combine-forward`, `mute-at-all`, `mute-red-envelope`, `group-role remove-user`, `group audit-join-validation` enum limits) are routed to the expected `im/chat` servers. They need permission rule or service-side confirmation.
- `chat group-mute-member`: old dynamic envelope mapped `--users` to `openDingTalkIds`; current helper also supports `uids`. The live error `uids is required` conflicts with the old schema, so this needs current backend `tools/list` schema confirmation before changing CLI mapping.
- PAT medium-risk nondeterminism is not caused by service discovery or CLI cache. CLI forwards the server risk result; policy needs server-side trace by `agentCode`, scope, corp/user, flowId/traceId.
- `chat group members remove` allowing the only owner to be removed, and `group update-icon` accepting a fake media id, are service-side validation gaps. CLI can add a guard only after product decides the exact rule.
