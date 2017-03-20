package cmd

import (
	"fmt"

	"github.com/akerl/speculate/utils"

	"github.com/spf13/cobra"
)

func envRunner(cmd *cobra.Command, args []string) error {
	rolename, err := utils.RoleNameParse(args)
	if err != nil {
		return err
	}
	role, err := utils.AssumeRole(rolename, cmd.Flags())
	if err != nil {
		return err
	}
	for _, line := range role.ToEnvVars() {
		fmt.Println(line)
	}
	return nil
}

var envCmd = &cobra.Command{
	Use:   "env ROLENAME",
	Short: "Print environment variables for an assumed role",
	RunE:  envRunner,
}

func init() {
	RootCmd.AddCommand(envCmd)
	utils.AddAssumeFlags(envCmd)
}
