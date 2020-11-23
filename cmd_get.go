package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

	rawFilePath := c.Args().Get(0)

	err := execGhqGet(rawFilePath)
	if err != nil {
		return err
	}

	filePath := trimFilePath(rawFilePath)
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

func trimFilePath(rawFilePath string) string {
	trimmedFilePath := strings.TrimPrefix(rawFilePath, urlPrefix)
	trimmedFilePath = strings.TrimSuffix(trimmedFilePath, gitSuffix)

	if !strings.HasPrefix(trimmedFilePath, hostPrefix) {
		trimmedFilePath = hostPrefix + "/" + trimmedFilePath
	}
	return trimmedFilePath
}

func execGhqGet(filePath string) error {
	ghqArgs := []string{
		"get",
		filePath,
	}
	fmt.Println("Fetching code...")

	cmd := exec.Command("ghq", ghqArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getGhqRoot() string {
	ghqRootArg := []string{"root"}
	ghqRoot, _ := exec.Command("ghq", ghqRootArg...).Output()
	rootPath := strings.TrimRight(string(ghqRoot), "\n")
	return rootPath
}

func getSccSummaryForHistory(fullPath string, filePath string) string {
	sccArgs := []string{
		"-f",
		"json",
		"-s",
		"lines",
		fullPath,
	}

	sccOutputForHistory, _ := exec.Command("scc", sccArgs...).Output()

	var summary []LanguageSummary
	json.Unmarshal(sccOutputForHistory, &summary)
	mainSummary := strings.TrimPrefix(filePath, hostPrefix+"/") + " " + summary[0].Name + " " + strconv.FormatInt(summary[0].Lines, 10)
	return mainSummary
}

type LanguageSummary struct {
	Name               string
	Bytes              int64
	Lines              int64
	Code               int64
	Comment            int64
	Blank              int64
	Complexity         int64
	Count              int64
	WeightedComplexity float64
	Files              interface{}
}
