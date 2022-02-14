package main

import (
	"github.com/eviltomorrow/blog-deploy-tool/pkg/cmd"
	"github.com/eviltomorrow/blog-deploy-tool/pkg/system"
)

var (
	GitSha      = ""
	GitTag      = ""
	GitBranch   = ""
	BuildTime   = ""
	MainVersion = "v1.0"
)

func setVersion() {
	system.CurrentVersion = MainVersion
	system.GitSha = GitSha
	system.GitTag = GitTag
	system.GitBranch = GitBranch
	system.BuildTime = BuildTime
}

func main() {
	setVersion()
	cmd.NewClient()
}
