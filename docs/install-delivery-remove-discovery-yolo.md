# Fork Delivery Install: remove-discovery + default-yolo

This fork-only preview combines:

- `feat/remove-discovery`: static endpoint mode and old open CLI compatibility.
- `feat/connect-yolo`: `dws dev connect` defaults to yolo/highest-permission mode.

Install:

```sh
curl -fsSL https://raw.githubusercontent.com/PeterGuy326/dingtalk-workspace-cli/delivery-remove-discovery-yolo-latest/scripts/install-delivery-remove-discovery-yolo.sh | sh
```

The installer downloads from:

```text
https://github.com/PeterGuy326/dingtalk-workspace-cli/releases/tag/delivery-remove-discovery-yolo-latest
```

Defaults:

- Binary: `~/.local/bin/dws`
- Skills: `dws skill setup --mode mono --target all --yes`
- Release assets: `dws-darwin-arm64.tar.gz`, `dws-darwin-amd64.tar.gz`, `dws-linux-arm64.tar.gz`, `dws-linux-amd64.tar.gz`, `dws-skills.zip`, `checksums.txt`

Optional:

```sh
# Skip skill setup.
curl -fsSL https://raw.githubusercontent.com/PeterGuy326/dingtalk-workspace-cli/delivery-remove-discovery-yolo-latest/scripts/install-delivery-remove-discovery-yolo.sh | DWS_NO_SKILLS=1 sh

# Install only Codex skill.
curl -fsSL https://raw.githubusercontent.com/PeterGuy326/dingtalk-workspace-cli/delivery-remove-discovery-yolo-latest/scripts/install-delivery-remove-discovery-yolo.sh | DWS_SKILL_TARGET=codex sh

# Use multi-skill layout.
curl -fsSL https://raw.githubusercontent.com/PeterGuy326/dingtalk-workspace-cli/delivery-remove-discovery-yolo-latest/scripts/install-delivery-remove-discovery-yolo.sh | DWS_SKILL_MODE=multi sh
```

Quick checks:

```sh
dws version
dws --help
dws dev connect --help
dws dev connect --channel codex --robot-client-id fake --robot-client-secret fake --dry-run --format json
dws dev connect --channel codex --robot-client-id fake --robot-client-secret fake --agent-permission-mode ask --dry-run --format json
```

Expected:

- `dws version` shows `MCP Static Endpoint Mode`.
- `dws --help` does not expose retired compatibility commands such as `cache` or `conference`.
- `dws dev connect --help` exposes `--yolo`, `--agent-permission-mode`, `--agent-approval-mode`; it does not expose `--agent-yolo`.
- Default dry-run is yolo/highest-permission mode. Passing `--agent-permission-mode ask` or `--agent-approval-mode ask` switches back to restricted/confirmation mode.
