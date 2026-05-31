package main

import (
	"os"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/editions"
	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/pkg/cli"
)

var (
	version   = "0.2.65"
	buildTime = "unknown"
	gitCommit = "unknown"
	buildMode = "dev"
)

func main() {
	editions.InstallEdition(buildMode, version)
	cli.SetVersion(version, buildTime, gitCommit)
	os.Exit(cli.Execute())
}
