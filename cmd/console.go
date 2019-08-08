package cmd

import (
	"fmt"

	"github.com/akerl/speculate/v2/creds"

	"github.com/spf13/cobra"
)

func consoleRunner(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()

	flagService, err := flags.GetString("service")
	if err != nil {
		return err
	}

	c, err := creds.NewFromEnv()
	if err != nil {
		return err
	}

	url, err := c.ToCustomConsoleURL(flagService)
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
	consoleCmd.Flags().String(
		"service", "", "Service path to access (defaults to AWS console homepage)")
}
