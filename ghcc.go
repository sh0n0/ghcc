package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	urlPrefix  = "https://"
	hostPrefix = "github.com"
	gitSuffix  = ".git"
)

func main() {
	app := &cli.App{
		Action: getCodeAndCount,
		Usage:  "ghq get and scc, and then rm if you wish",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "temporary", Aliases: []string{"t"}, Usage: "remove code after counting"},
		},
		Commands: []*cli.Command{
			{
				Name:   "ls",
				Usage:  "list history of scc",
				Action: ls,
			},
		},
	}
	app.Run(os.Args)
}

func getCodeAndCount(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("Just one argument is required")
		return nil
	}

	filePath := c.Args().Get(0)
	ghqArgs := []string{
		"get",
		filePath,
	}
	ghqOut, _ := exec.Command("ghq", ghqArgs...).CombinedOutput()
	fmt.Println(string(ghqOut))

	filePath = strings.TrimPrefix(filePath, urlPrefix)
	filePath = strings.TrimSuffix(filePath, gitSuffix)

	if !strings.HasPrefix(filePath, hostPrefix) {
		filePath = hostPrefix + "/" + filePath
	}

	ghqRootArg := []string{"root"}
	ghqRoot, _ := exec.Command("ghq", ghqRootArg...).Output()
	rootPath := strings.TrimRight(string(ghqRoot), "\n")

	fullPath := rootPath + "/" + filePath
	out, _ := exec.Command("scc", fullPath).Output()
	fmt.Println(string(out))

	isTemporary := c.Bool("temporary")
	if isTemporary {
		os.RemoveAll(fullPath)
		fmt.Println("Finished removing " + fullPath)
		return nil
	}

interactive_loop:
	for {
		stdin := bufio.NewScanner(os.Stdin)
		fmt.Print("Remove source code? (y/n): ")
		stdin.Scan()
		switch stdin.Text() {
		case "y":
			os.RemoveAll(fullPath)
			fmt.Println("Finished removing " + fullPath)
			break interactive_loop
		case "n":
			break interactive_loop
		default:
			fmt.Println("Please enter y or n")
		}
	}
	return nil
}

func ls(c *cli.Context) error {
	// TODO
	return nil
}
