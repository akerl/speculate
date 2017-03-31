package cmd

import (
	"fmt"

	"github.com/akerl/speculate/utils"

	"github.com/spf13/cobra"
)

func consoleRunner(cmd *cobra.Command, args []string) error {
	role, err := utils.NewRoleFromEnv()
	if err != nil {
		return err
	}
	url, err := role.ToConsoleURL()
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
