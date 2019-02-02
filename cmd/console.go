package cmd

import (
	"fmt"

	"github.com/akerl/speculate/creds"

	"github.com/spf13/cobra"
)

func consoleRunner(_ *cobra.Command, _ []string) error {
	c, err := creds.NewFromEnv()
	if err != nil {
		return err
	}
	url, err := c.ToConsoleURL()
	if err != nil {
		return err
	}

	fmt.Println(url)
	return nil
}

var consoleCmd = &cobra.Command{
	Use:   "console ROLENAME",
	Short: "Open current role in browser",
	RunE:  consoleRunner,
}

func init() {
	rootCmd.AddCommand(consoleCmd)
}
