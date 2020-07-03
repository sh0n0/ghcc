package main

import "github.com/urfave/cli/v2"

var commands = []*cli.Command{
	{
		Name:   "get",
		Usage:  "ghq get and scc, and then rm if you wish",
		Action: getCodeAndCount,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "temporary", Aliases: []string{"t"}, Usage: "remove code after counting"},
		},
	},
	{
		Name:   "ls",
		Usage:  "list history of scc",
		Action: ls,
	},
	{
		Name:   "clear",
		Usage:  "clear all history",
		Action: clear,
	},
}
