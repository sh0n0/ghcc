package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	urlPrefix       = "https://"
	hostPrefix      = "github.com"
	gitSuffix       = ".git"
	historyFileName = ".ghcc_history"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
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
		},
	}
	app.Run(os.Args)
}

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

func ls(c *cli.Context) error {
	file, err := os.Open(getHistoryPath())
	if os.IsNotExist(err) {
		fmt.Println("No history")
		return nil
	}

	if err != nil {
		return err
	}

	defer file.Close()

	history, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fmt.Println(string(history))
	return nil
}

func getHistoryPath() string {
	homePath, _ := os.UserHomeDir()
	return filepath.Join(homePath, historyFileName)
}

// Add new result to history
// If a history file doesn't exist, create new one
func writeResultToHistory(result string) {
	file, _ := os.OpenFile(getHistoryPath(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer file.Close()

	fmt.Fprintln(file, result)
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
	sortLaunguageSumamary(summary)
	mainSummary := strings.TrimPrefix(filePath, hostPrefix) + " " + summary[0].Name + " " + strconv.FormatInt(summary[0].Lines, 10)
	return mainSummary
}

// this will be removed when the next version of scc releases
func sortLaunguageSumamary(summary []LanguageSummary) {
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].Lines > summary[j].Lines
	})
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
