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

func getCodeAndCount(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("Just one argument is required")
		return nil
	}

	rowFilePath := c.Args().Get(0)

	execGhqGet(rowFilePath)

	filePath := trimFilePath(rowFilePath)
	fullPath := getGhqRoot() + "/" + filePath

	showSccResult(fullPath)

	mainSummary := getSccSummaryForHistory(fullPath, filePath)
	writeResultToHistory(mainSummary)

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

func showSccResult(fullPath string) {
	sccArgs := []string{
		"-s",
		"lines",
		fullPath,
	}
	sccOutput, _ := exec.Command("scc", sccArgs...).Output()
	fmt.Println(string(sccOutput))
}

func trimFilePath(rowFilePath string) string {
	rowFilePath = strings.TrimPrefix(rowFilePath, urlPrefix)
	rowFilePath = strings.TrimSuffix(rowFilePath, gitSuffix)

	if !strings.HasPrefix(rowFilePath, hostPrefix) {
		rowFilePath = hostPrefix + "/" + rowFilePath
	}
	return rowFilePath
}

func execGhqGet(filePath string) {
	ghqArgs := []string{
		"get",
		filePath,
	}
	fmt.Println("Fetching code...")
	ghqOut, _ := exec.Command("ghq", ghqArgs...).CombinedOutput()
	fmt.Println(string(ghqOut))
}

func getGhqRoot() string {
	ghqRootArg := []string{"root"}
	ghqRoot, _ := exec.Command("ghq", ghqRootArg...).Output()
	rootPath := strings.TrimRight(string(ghqRoot), "\n")
	return rootPath
}
