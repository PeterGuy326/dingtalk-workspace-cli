//go:build !wukong

package editions

import "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/edition"

// hooks returns nil for the open-source edition: edition.Override(nil) is a
// no-op, so the zero-value open-source defaults (pkg/edition.defaultHooks)
// stay active. buildMode is accepted for signature parity but unused — the
// open edition behaves identically in dev/real.
func hooks(buildMode, version string) *edition.Hooks {
	return nil
}
