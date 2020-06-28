package main

import (
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
	if c.Args().Len() != 1 {
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

	sccArgs := rootPath + "/" + filePath
	out, _ := exec.Command("scc", sccArgs).Output()
	fmt.Println(string(out))
	return nil
}

func ls(c *cli.Context) error {
	out, _ := exec.Command("pwd").Output()
	fmt.Println(string(out))
	return nil
}
