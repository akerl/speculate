package cmd

import (
	"fmt"

	"github.com/akerl/speculate/executors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func envRunner(cmd *cobra.Command, args []string) error {
	var exec executors.Executor
	var err error

	switch len(args) {
	case 0:
		exec = &executors.Signin{}
	case 1:
		exec = &executors.Assumption{}
		if err := exec.SetRoleName(args[0]); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Too many args provided. Check --help for more info")
	}

	flags := cmd.Flags()
	if err := parseFlags(exec, flags); err != nil {
		return err
	}

	creds, err := exec.Execute()
	if err != nil {
		return err
	}
	for _, line := range creds.ToEnvVars() {
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
	envCmd.Flags().StringP("account", "a", "", "Account ID to assume role on (defaults to source account")
	envCmd.Flags().StringP("session", "s", "", "Set session name for assumed role (defaults to origin user name)")
	envCmd.Flags().Int64P("lifetime", "l", 3600, "Set lifetime of credentials in seconds (defaults to 3600 seconds / 1 hour, min 900, max 3600)")
	envCmd.Flags().String("policy", "", "Set a IAM policy in JSON for the assumed credentials")
	envCmd.Flags().BoolP("mfa", "m", false, "Use MFA when assuming role")
	envCmd.Flags().String("mfacode", "", "Code to use for MFA")
}

func parseFlags(exec executors.Executor, flags *pflag.FlagSet) error {
	if val, err := flags.GetString("account"); err != nil {
		return err
	} else if val != "" {
		if err := exec.SetAccountID(val); err != nil {
			return err
		}
	}

	if val, err := flags.GetString("session"); err != nil {
		return err
	} else if val != "" {
		if err := exec.SetSessionName(val); err != nil {
			return err
		}
	}

	if val, err := flags.GetString("policy"); err != nil {
		return err
	} else if val != "" {
		if err := exec.SetPolicy(val); err != nil {
			return err
		}
	}

	if val, err := flags.GetInt64("lifetime"); err != nil {
		return err
	} else if val != 0 {
		if err := exec.SetLifetime(val); err != nil {
			return err
		}
	}

	if val, err := flags.GetBool("mfa"); err != nil {
		return err
	} else if err := exec.SetMfa(val); err != nil {
		return err
	}

	if val, err := flags.GetString("mfacode"); err != nil {
		return err
	} else if val != "" {
		if err := exec.SetMfaCode(val); err != nil {
			return err
		}
	}

	return nil
}
