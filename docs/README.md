# Documentation

This directory contains durable documentation that is not already available from the CLI itself.

Command names, flags, and subcommands are intentionally not copied into a static command catalog. Use `dws --help` and `dws <command> --help`; those commands are built from the same command tree that executes requests and are the source of truth.

## User guides

- [Robot quickstart](./guides/robot-quickstart.md) — create a DingTalk robot and connect it to a local AI agent.
- [Run the connector as a service](./guides/connect-daemon.md) — keep an existing connector running in the background or across reboots.
- [CLI reference](./reference.md) — stable global output, error, schema, and completion contracts.

## Architecture and integrations

- [Architecture](./architecture.md) — repository structure and request flow.
- [Agent identification](./integrations/agent-identification.md) — wire-level analytics headers and their trust boundaries.

## Maintainers

- [Release guide](./maintainers/releasing.md) — prerelease and stable release procedure.
