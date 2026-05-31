// Package editions performs compile-time edition selection. The active edition
// is chosen by build tag: building with `-tags wukong` links the wukong overlay
// (selector_wukong.go); building without it links the open-source default
// (selector_open.go). main calls InstallEdition exactly once before Execute.
package editions

import "github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/edition"

// InstallEdition installs the edition-specific hooks selected at link time.
// The build-tagged hooks() returns the wukong overlay hooks under `-tags
// wukong`, or nil otherwise. edition.Override is a no-op on nil, so the open
// build keeps the zero-value open-source defaults.
func InstallEdition(buildMode, version string) {
	edition.Override(hooks(buildMode, version))
}
