package main

import (
	"os"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/cli"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/edition"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/wukong"
)

var (
	version   = "0.2.30"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	edition.Override(wukong.NewHooks())
	cli.SetVersion(version, buildTime, gitCommit)
	os.Exit(cli.Execute())
}
