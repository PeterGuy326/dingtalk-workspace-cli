#!/bin/sh
# Copyright 2026 Alibaba Group
# Licensed under the Apache License, Version 2.0
#
# Preview installer for the feat/remove-discovery branch from the PeterGuy326 fork.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/PeterGuy326/dingtalk-workspace-cli/feat-remove-discovery-latest/scripts/install-feat-remove-discovery.sh | sh
#
# Optional environment variables:
#   DWS_REPO            owner/repo that hosts the release assets (default: PeterGuy326/dingtalk-workspace-cli)
#   DWS_VERSION         release tag to install (default: feat-remove-discovery-latest)
#   DWS_INSTALL_DIR     install directory (default: ~/.local/bin)
#   DWS_INSTALL_NAME    installed binary name (default: dws)
#   DWS_NO_SKILLS       set 1 to skip skill setup
#   DWS_SKILL_MODE      mono | multi (default: mono)
#   DWS_SKILL_TARGET    all | codex | claude | cursor | ... (default: all)

set -eu

REPO="${DWS_REPO:-PeterGuy326/dingtalk-workspace-cli}"
VERSION="${DWS_VERSION:-feat-remove-discovery-latest}"
INSTALL_DIR="${DWS_INSTALL_DIR:-$HOME/.local/bin}"
INSTALL_NAME="${DWS_INSTALL_NAME:-dws}"
NO_SKILLS="${DWS_NO_SKILLS:-0}"
SKILL_MODE="${DWS_SKILL_MODE:-mono}"
SKILL_TARGET="${DWS_SKILL_TARGET:-all}"

say() {
  printf '  %s\n' "$@"
}

err() {
  printf '  ERROR: %s\n' "$@" >&2
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1
}

download() {
  url="$1"
  dest="$2"
  if need_cmd curl; then
    curl -fsSL "$url" -o "$dest"
    return 0
  fi
  if need_cmd wget; then
    wget -qO "$dest" "$url"
    return 0
  fi
  err "curl or wget is required"
}

detect_os() {
  case "$(uname -s)" in
    Darwin*) printf 'darwin\n' ;;
    Linux*) printf 'linux\n' ;;
    *) err "unsupported OS: $(uname -s)" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    arm64|aarch64) printf 'arm64\n' ;;
    x86_64|amd64) printf 'amd64\n' ;;
    *) err "unsupported architecture: $(uname -m)" ;;
  esac
}

verify_checksum() {
  file="$1"
  checksums="$2"
  name="$(basename "$file")"
  expected="$(grep "  ${name}\$" "$checksums" | awk '{print $1}' | head -1 || true)"
  if [ -z "$expected" ]; then
    say "WARN: ${name} not found in checksums.txt; skipping checksum verification"
    return 0
  fi

  if need_cmd shasum; then
    actual="$(shasum -a 256 "$file" | awk '{print $1}')"
  elif need_cmd sha256sum; then
    actual="$(sha256sum "$file" | awk '{print $1}')"
  else
    say "WARN: shasum/sha256sum not found; skipping checksum verification"
    return 0
  fi

  [ "$actual" = "$expected" ] || err "checksum mismatch for ${name}"
}

extract_binary() {
  archive="$1"
  dest="$2"
  mkdir -p "$dest"
  tar -xzf "$archive" -C "$dest"
  if [ -f "$dest/dws" ]; then
    printf '%s\n' "$dest/dws"
    return 0
  fi
  found="$(find "$dest" -type f -name dws | head -1)"
  [ -n "$found" ] || err "archive does not contain dws binary"
  printf '%s\n' "$found"
}

case "$SKILL_MODE" in
  mono|multi) ;;
  *) err "DWS_SKILL_MODE must be mono or multi" ;;
esac

OS="$(detect_os)"
ARCH="$(detect_arch)"
ASSET="dws-${OS}-${ARCH}.tar.gz"
BASE_URL="https://github.com/${REPO}/releases/download/${VERSION}"
TMPDIR="$(mktemp -d 2>/dev/null || mktemp -d -t dws-install)"

cleanup() {
  rm -rf "$TMPDIR"
}
trap cleanup EXIT INT TERM

printf '\n'
say "DWS feat/remove-discovery installer"
say "repo:    ${REPO}"
say "version: ${VERSION}"
say "target:  ${OS}/${ARCH}"
printf '\n'

download "${BASE_URL}/${ASSET}" "$TMPDIR/$ASSET"
download "${BASE_URL}/checksums.txt" "$TMPDIR/checksums.txt"
verify_checksum "$TMPDIR/$ASSET" "$TMPDIR/checksums.txt"

bin_path="$(extract_binary "$TMPDIR/$ASSET" "$TMPDIR/extract")"
mkdir -p "$INSTALL_DIR"
cp "$bin_path" "$INSTALL_DIR/$INSTALL_NAME"
chmod +x "$INSTALL_DIR/$INSTALL_NAME"

say "Installed: $INSTALL_DIR/$INSTALL_NAME"

if [ "$NO_SKILLS" != "1" ]; then
  say "Installing embedded skills: mode=${SKILL_MODE}, target=${SKILL_TARGET}"
  "$INSTALL_DIR/$INSTALL_NAME" skill setup --mode "$SKILL_MODE" --target "$SKILL_TARGET" --yes
else
  say "Skipped skill setup because DWS_NO_SKILLS=1"
fi

printf '\n'
"$INSTALL_DIR/$INSTALL_NAME" version
printf '\n'
say "Done. If ${INSTALL_DIR} is not on PATH, add it or run: ${INSTALL_DIR}/${INSTALL_NAME}"
