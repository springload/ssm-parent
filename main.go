package main

import "github.com/springload/ssm-parent/cmd"

var version = "dev"

func main() {
	cmd.Execute(version)
}
