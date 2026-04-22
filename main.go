package main

import (
	"os"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/cli"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/edition"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/wukong"
)

var (
	version   = "0.2.50.1"
	buildTime = "unknown"
	gitCommit = "unknown"
	buildMode = "dev"
)

func main() {
	edition.Override(wukong.NewHooks(buildMode))
	cli.SetVersion(version, buildTime, gitCommit)
	os.Exit(cli.Execute())
}
