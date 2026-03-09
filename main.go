package main

import "github.com/bithostio/bh/cmd"

var version = "dev"

func main() {
	cmd.Version = version
	cmd.Execute()
}
