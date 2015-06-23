package main

import "gopkg.in/acd.v0/cli"

func main() {
	app := cli.New()
	app.RunAndExitOnError()
}
