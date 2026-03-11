package main

import "github.com/optimode/ipinfo/cmd"

// Set via ldflags at build time.
var (
	version   = "dev"
	gitCommit = "unknown"
	buildTime = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, gitCommit, buildTime)
	cmd.Execute()
}
