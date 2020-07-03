package main

import (
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

const historyFileName = ".ghcc_history"

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
