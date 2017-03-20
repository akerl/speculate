package cmd

import (
	"fmt"

	"github.com/akerl/speculate/utils"

	"github.com/spf13/cobra"
)

var logout_url = "https://console.aws.amazon.com/ec2/logout!doLogout"

func consoleRunner(cmd *cobra.Command, args []string) error {
	role, err := utils.NewRoleFromEnv()
	if err != nil {
		return err
	}
	url, err := role.ToConsoleUrl()
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
	RootCmd.AddCommand(consoleCmd)
	consoleCmd.Flags().BoolP("logout", "l", false, "Open logout page before opening new session URL")
}
