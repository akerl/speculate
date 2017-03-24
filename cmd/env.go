package cmd

import (
	"fmt"

	"github.com/akerl/speculate/utils"

	"github.com/spf13/cobra"
)

func envRunner(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	accountID, err := flags.GetString("account")
	if err != nil {
		return err
	}
	sessionName, err := flags.GetString("session")
	if err != nil {
		return err
	}
	useMfa, err := flags.GetBool("mfa")
	if err != nil {
		return err
	}
	mfaCode, err := flags.GetString("mfacode")
	if err != nil {
		return err
	}
	roleName, err := utils.RoleNameParse(args)
	if err != nil {
		return err
	}
	role, err := utils.AssumeRole(roleName, accountID, sessionName, useMfa, mfaCode)
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
