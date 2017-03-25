package cmd

import (
	"fmt"

	"github.com/akerl/speculate/utils"

	"github.com/spf13/cobra"
)

func envRunner(cmd *cobra.Command, args []string) error {
	assumption := utils.Assumption{}
	if len(args) < 1 {
		return fmt.Errorf("No role name provided")
	}
	assumption.RoleName = args[0]
	role, err := assumption.AssumeRole()
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
	rootCmd.AddCommand(envCmd)
	utils.AddAssumeFlags(envCmd)
}
