package main

import "github.com/wim-vdw/wimkube/cmd"

var version = "developer build"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
