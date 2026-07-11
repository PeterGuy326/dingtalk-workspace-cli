# DingTalk AI Robot Quickstart

Create a DingTalk group robot and connect it to a local AI CLI in four steps. The connector supports built-in channels such as Codex, Claude Code, Qoder, Gemini, and opencode, plus custom one-shot commands.

## 1. Install and sign in

macOS or Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/DingTalk-Real-AI/dingtalk-workspace-cli/main/scripts/install-devapp.sh | sh
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/DingTalk-Real-AI/dingtalk-workspace-cli/main/scripts/install-devapp.ps1 | iex
```

Open a new terminal if the installer changed `PATH`, then verify the installation and sign in:

```bash
dws version
dws auth login
```

The convenience installer installs the current stable `dws` binary and the `dingtalk-dev` Agent Skill. The standard repository installer also includes `dws dev`.

## 2. Create the robot

Robot creation is asynchronous. Preview the request first:

```bash
dws dev app robot submit \
  --name "My Assistant" \
  --robot-name "Team Helper" \
  --desc "Answer questions in the group" \
  --dry-run \
  --format json
```

After reviewing the target and values, execute it:

```bash
dws dev app robot submit \
  --name "My Assistant" \
  --robot-name "Team Helper" \
  --desc "Answer questions in the group" \
  --yes \
  --format json
```

Save the returned `taskId`, then poll until `status` is `SUCCESS`:

```bash
dws dev app robot result --task-id TASK_ID --format json
```

Save the returned `unifiedAppId`. Do not copy `clientSecret` into a shell command or document.

## 3. Connect the robot to a local agent

```bash
dws dev connect --channel auto --unified-app-id UNIFIED_APP_ID
```

`--channel auto` selects an installed supported AI CLI. Using `--unified-app-id` lets `dws` resolve credentials at runtime, so the secret is not exposed through the process list or shell history.

The command runs in the foreground. Keep the terminal open during this first test, then send a message to confirm the local agent can answer. For a background or boot-time service, follow [Run the connector as a service](./connect-daemon.md).

For a custom one-shot AI command, use:

```bash
dws dev connect \
  --agent-cmd "your-ai-command --prompt" \
  --unified-app-id UNIFIED_APP_ID
```

The incoming question is appended as the final argument and the command's stdout becomes the reply. Run `dws dev connect --help` for optional access controls, knowledge sources, model selection, and approval settings.

## 4. Add the robot to a group

In DingTalk, open the target group and choose:

```text
Group Settings → Robots → Add Robot → Enterprise Robots → Team Helper
```

Add the robot, then @-mention it in the group. The connector terminal should show the incoming request and reply.

## Troubleshooting the quickstart

- `dws dev app robot submit` requires Open Platform developer permission. Ask an enterprise administrator to add the current account as a developer if the command rejects the user role.
- A `WAITING` task is not a failure; poll `robot result` again after the returned interval. Use credentials only after `SUCCESS`.
- If `--channel auto` finds no agent, install and authenticate a supported AI CLI or provide `--agent-cmd`.
- If the group robot stops replying, confirm the foreground connector is still running. For a daemon, use `dws dev connect status --unified-app-id UNIFIED_APP_ID`.
- `dws dev connect` does not create an approval request. Robot creation and application publication are the operations that may involve platform approval.
