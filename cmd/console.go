package cmd

import (
	"github.com/akerl/speculate/utils"

	"github.com/spf13/cobra"
)

func consoleRunner(cmd *cobra.Command, args []string) {
}

var consoleCmd = &cobra.Command{
	Use:   "console ROLENAME",
	Short: "Open assumed role in browser",
	Run:   consoleRunner,
}

func init() {
	RootCmd.AddCommand(consoleCmd)
	utils.AddAssumeFlags(consoleCmd)
	consoleCmd.Flags().BoolP("logout", "l", false, "Open logout page before opening new session URL")
}
