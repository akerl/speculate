package cmd

import (
	"fmt"

	"github.com/akerl/speculate/utils"

	"github.com/spf13/cobra"
)

func envRunner(cmd *cobra.Command, args []string) {
	role := utils.RoleParse(args)
	assumption := utils.AssumeRole(role, cmd.Flags())
	fmt.Printf("Access Key: %s\n", assumption.Access_key)
}

var envCmd = &cobra.Command{
	Use:   "env ROLENAME",
	Short: "Print environment variables for an assumed role",
	Run:   envRunner,
}

func init() {
	RootCmd.AddCommand(envCmd)
	utils.AddAssumeFlags(envCmd)
}
