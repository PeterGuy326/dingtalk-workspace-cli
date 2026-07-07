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

Release binary checked:

```json
{
  "architecture": "MCP Static Endpoint Mode",
  "commit": "1b823514",
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
| `dws cache refresh` compatibility stub | Pass |
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

## Skill Setup Verification

Mono mode:

```text
HOME=<tmp> dws skill setup --mode mono --target codex --yes
installed: <tmp>/.codex/skills/dws/SKILL.md
file count: 133
```

Multi mode:

```text
HOME=<tmp> dws skill setup --mode multi --target codex --yes -s dev -s chat
installed:
  <tmp>/.codex/skills/dws-shared/SKILL.md
  <tmp>/.codex/skills/dingtalk-chat/SKILL.md
  <tmp>/.codex/skills/dingtalk-dev/SKILL.md
file count: 21
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

Summary:

| Tree | Visible commands | Help errors |
| --- | ---: | ---: |
| Open delivery | 779 | 0 |
| Wukong local build | 973 | 0 |

Skill/help drift notes:

- No install drift: mono and multi skill setup both install successfully in an isolated HOME.
- No command-help reachability drift: every visible open delivery command path has a working `dws <path> --help`.
- Content drift found in mono skill coverage: open help exposes `live` and `pat`, but mono skill does not currently include product reference pages for them. Multi mode has `dingtalk-live`; `pat` is still not covered by product skill docs.

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
- `dws-wukong origin/develop` at `2a7c70ff` currently does not compile directly against the remove-discovery core because it still references removed `pkg/cmdutil` overlay APIs: `IsEnvelopeSourced` and `MergeHardcodedLeaves`.
- The wukong command tree snapshot was generated from the local wukong workspace at `f48ae79d` with the local compatibility bridge changes present in that workspace.
