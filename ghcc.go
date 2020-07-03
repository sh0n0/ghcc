package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "ghcc"
	app.Commands = commands
	app.Run(os.Args)
}
