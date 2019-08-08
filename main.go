package main

import (
	"os"

	"github.com/akerl/speculate/cmd"
	"github.com/akerl/speculate/helpers"
)

func main() {
	if err := cmd.Execute(); err != nil {
		helpers.PrintAwsError(err)
		os.Exit(1)
	}
}
