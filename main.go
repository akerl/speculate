package main

import (
	"os"

	"github.com/akerl/speculate/v2/cmd"
	"github.com/akerl/speculate/v2/helpers"
)

func main() {
	if err := cmd.Execute(); err != nil {
		helpers.PrintAwsError(err)
		os.Exit(1)
	}
}
