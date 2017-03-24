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
	lifetime, err := flags.GetInt64("lifetime")
	if err != nil {
		return err
	}
	if lifetime < 900 || lifetime > 3600 {
		return fmt.Errorf("Lifetime must be between 900 and 3600: %d", lifetime)
	}
	role, err := utils.AssumeRole(roleName, accountID, sessionName, useMfa, mfaCode, lifetime)
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
