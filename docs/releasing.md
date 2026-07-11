# Release Guide (Prerelease / Stable)

All releases follow one path: the local script seals the release candidate, verifies it, and pushes an annotated tag; GitHub Actions builds and publishes the final artifacts. Do not run `goreleaser release` directly, and never create, replace, or move tags by hand.

Before releasing, complete the platform governance prerequisites: immutable releases must be enabled in the target GitHub repository, `main` must require `CI Gate`, and the release machine must have `gh` installed and authenticated. Before it creates a tag, the local script calls the API to check immutable releases, the `CI Gate` status for the current SHA, and any in-progress Release workflow. A `v*` tag ruleset must still be configured in advance by a repository administrator and confirmed by the operator.

## One entry point for everyday use

After installing the release Skill, run:

```bash
dws-release
```

With no arguments, it starts in guided mode. The equivalent in-repository command is `./scripts/release/dws-release.sh`. The first use requires a one-time configuration of the production publishing remote. The command stores both the remote name and its canonical repository identity in the current Git repository:

```bash
dws-release config --remote origin
```

After that, the command chooses the right step from repository state: if the exact CHANGELOG section is missing, it only creates a template and stops; after the section is completed, committed, and merged to `main`, rerunning the same command safely fast-forwards local `main` and completes the full preflight. It rejects a configured remote if it is later retargeted to a different repository. Publishing requires an explicit `--publish`, and the underlying command still requires final version confirmation.

## Release model

```text
Release candidate on main + beta CHANGELOG
  → vX.Y.Z-beta.N (prerelease validation)
  → only the stable CHANGELOG may be added; source code must not change
  → vX.Y.Z (stable release)
```

The stable release must explicitly name the beta that was validated. The script compares them and refuses the stable release if any file other than `CHANGELOG.md` has changed. The code and command tree tested in prerelease are therefore exactly the code and command tree delivered in the stable release.

## Prerelease

Run the common entry point:

```bash
dws-release v1.2.3-beta.1
```

If the CHANGELOG section does not exist, this command only creates a template and stops. Complete it, remove every `TODO`, commit it, and merge it to `main` through a PR. Then run the exact same command again; it performs the complete preflight:

```bash
dws-release v1.2.3-beta.1
```

The preflight includes tests, policy checks, command-tree compatibility against the previous stable release, packaging for all target platforms, npm installation verification, and Homebrew installation verification on macOS. When it passes, publish with:

```bash
dws-release v1.2.3-beta.1 --publish
```

After all preflight steps pass, the command asks again for the complete version. The common entry point does not support bypassing this confirmation.

## Stable release

After validating the beta, run the stable entry point:

```bash
dws-release v1.2.3 --from-beta v1.2.3-beta.1
```

The first run only creates the stable CHANGELOG template and stops. Complete it, remove `TODO`, commit it, and merge it to `main` through a PR. Rerun the same command for the full preflight, then add `--publish` after confirmation:

```bash
dws-release v1.2.3 --from-beta v1.2.3-beta.1
dws-release v1.2.3 --from-beta v1.2.3-beta.1 --publish
```

`FROM_BETA` is never inferred automatically. It is written as `From-Beta` metadata on the stable annotated tag, then read and verified again by CI.

## CHANGELOG contract

Each tag requires one exact section that is unique, non-empty, and contains no `TODO` or `TBD`:

```markdown
## [1.2.3-beta.1] - 2026-07-11

### Changed

- User-visible changes validated by this beta.
```

For stable releases, use `## [1.2.3] - YYYY-MM-DD`. This section becomes the GitHub Release notes verbatim.

## CI/CD guarantees

- Only `vX.Y.Z-beta.N` and `vX.Y.Z` are accepted, and a new version must be greater than the preceding stable version. That preceding stable version must have both a public, non-draft GitHub Release and a successful Release workflow for the same tag and commit. An orphaned version with only a tag, but no successful delivery, blocks later releases until it is rerun and completed.
- Tags must be annotated. Before pushing, the local script rechecks that `HEAD` exactly matches remote `main`. CI allows `main` to advance afterward, but requires the sealed commit to remain in `main` history.
- Both ordinary CI and the release preflight compare the full command tree with the latest delivered stable release. If that baseline changes during a long preflight, the comparison runs again against the new baseline.
- GoReleaser only builds. Final GitHub Release artifacts are uploaded only after Darwin re-signing, checksum regeneration, and npm installation verification pass.
- Each of the six platform archives is unpacked and checked for its embedded binary version. The public asset set, checksum set, and npm tarball integrity must match exactly. The npm tarball is always packed with npm `10.9.2`, avoiding byte drift if a workflow is rerun on a runner with a different bundled npm version.
- Stable releases publish to npm `latest` and update OSS `latest.txt` plus the shared install script. Prereleases publish to npm `beta` and update only OSS `beta.txt`; they never overwrite the stable entry point.
- The Release workflow has a serial publication queue that can hold at most 100 pending runs. The local entry point still requires the preceding Release workflow to finish before it can seal the next tag.
- If a local tag push fails, the newly created local tag is removed. Once a tag push succeeds, release delivery belongs exclusively to CI: do not move the tag or reuse its version.

An npm repair is allowed only through the Release workflow's `repair_npm_version`, triggered from the default branch. It supports only public immutable releases produced successfully by this pipeline after immutable releases were enabled: the target must be an annotated tag in `main` history, and the `Build immutable GitHub Release` job for the same commit must have succeeded. This immutable artifact boundary remains valid even if later npm distribution fails. Repair reconstructs the package from the target commit's npm template, verifies all platform assets and binary versions, then publishes to the isolated `backfill` dist-tag. It does not roll back `latest` or `beta`. Historical mutable releases are intentionally excluded from automatic repair so replaceable assets cannot enter npm.

If OSS or Gitee distribution fails, rerun the failed `Publish npm and mirrors` job for that tag. Each step reuses the immutable GitHub artifacts and preserves channel monotonicity. The separate Gitee release workflow and local direct-publish script are retired; neither may bypass the publication queue or overwrite a mirror with different bytes from a rebuild.

OSS `latest.txt` and `beta.txt` are currently mirror-channel metadata. The repository installer still resolves versions from GitHub/Gitee, so an OSS pointer must not be treated as an enabled installation channel.

Homebrew is currently only a local preflight and manual-formula path: the preflight performs a real install for the current macOS architecture, but the Release workflow does not publish a tap. A single-host formula generated in CI must not be treated as a formal two-architecture Darwin delivery. The automatic stable delivery scope is GitHub Release, npm, OSS, and the Gitee fallback when explicitly enabled. Publishing a dual-architecture Homebrew tap requires a separate request.

## Platform governance prerequisites

Repository administrators must configure two platform rules that the script cannot replace:

- `main` must require the exact `CI Gate`; the tag workflow also verifies through the Checks API that this sealed SHA passed it.
- Immutable releases must be enabled. They protect only releases created after activation, so enable them before the first release through this pipeline. Add a `v*` tag ruleset to restrict creation permission and protect the short interval before a release is published.

When immutable releases or `CI Gate` is missing, the release script refuses to create a tag. A tag ruleset can be inherited from the organization, so the script does not infer its effective scope. Administrator confirmation remains required; a script convention cannot replace platform enforcement.
