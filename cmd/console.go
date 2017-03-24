package cmd

import (
	"fmt"

	"github.com/akerl/speculate/utils"

	"github.com/spf13/cobra"
)

func consoleRunner(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	lifetime, err := flags.GetInt("lifetime")
	if err != nil {
		return err
	}
	if lifetime > 43200 {
		return fmt.Errorf("Lifetime larger than max of 43200: %d", lifetime)
	}
	role, err := utils.NewRoleFromEnv()
	if err != nil {
		return err
	}
	url, err := role.ToConsoleURL(lifetime)
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
	consoleCmd.Flags().IntP("lifetime", "l", 3600, "Set lifetime of session in seconds (defaults to 3600 seconds / 1 hour, max 43200)")
}
