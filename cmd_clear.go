package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func clear(c *cli.Context) error {
	err := os.Remove(getHistoryPath())
	if err != nil {
		return err
	}
	fmt.Println("clear all history")
	return nil
}
