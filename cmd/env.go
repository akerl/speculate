package cmd

import (
	"fmt"

	"github.com/akerl/speculate/v2/creds"

	"github.com/spf13/cobra"
)

// revive:disable-next-line:cyclomatic
func envRunner(cmd *cobra.Command, args []string) error {
	var err error
	c := creds.Creds{}
	flags := cmd.Flags()

	accountID, err := flags.GetString("account")
	if err != nil {
		return err
	}

	session, err := flags.GetString("session")
	if err != nil {
		return err
	}

	policy, err := flags.GetString("policy")
	if err != nil {
		return err
	}

	lifetime, err := flags.GetInt64("lifetime")
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

	switch len(args) {
	case 0:
		if accountID != "" {
			return fmt.Errorf("accountID is invalid for session tokens")
		}
		if session != "" {
			return fmt.Errorf("session is invalid for session tokens")
		}
		if policy != "" {
			return fmt.Errorf("policy is invalid for session tokens")
		}
		c, err = c.GetSessionToken(creds.GetSessionTokenOptions{
			Lifetime: lifetime,
			UseMfa:   useMfa,
			MfaCode:  mfaCode,
		})
	case 1:
		c, err = c.AssumeRole(creds.AssumeRoleOptions{
			RoleName:    args[0],
			AccountID:   accountID,
			SessionName: session,
			Policy:      policy,
			Lifetime:    lifetime,
			UseMfa:      useMfa,
			MfaCode:     mfaCode,
		})
	default:
		return fmt.Errorf("too many args provided. Check --help for more info")
	}

	if err != nil {
		return err
	}

	for _, line := range c.ToEnvVars() {
		fmt.Println(line)
	}
	return nil
}

var envCmd = &cobra.Command{
	Use:   "env [ROLENAME]",
	Short: "Generate temporary credentials by assuming a role or requesting a session token",
	RunE:  envRunner,
}

func init() {
	rootCmd.AddCommand(envCmd)
	//revive:disable:line-length-limit
	envCmd.Flags().StringP("account", "a", "", "Account ID to assume role on (defaults to source account)")
	envCmd.Flags().StringP("session", "s", "", "Set session name for assumed role (defaults to origin user name)")
	envCmd.Flags().Int64P(
		"lifetime", "l",
		0,
		fmt.Sprintf(
			"Set lifetime of credentials in seconds. For SessionToken, must be between %d and %d (default %d). For AssumeRole, must be between %d and %d (default %d)",
			creds.SessionTokenLifetimeLimits.Min,
			creds.SessionTokenLifetimeLimits.Max,
			creds.SessionTokenLifetimeLimits.Default,
			creds.AssumeRoleLifetimeLimits.Min,
			creds.AssumeRoleLifetimeLimits.Max,
			creds.AssumeRoleLifetimeLimits.Default,
		),
	)
	envCmd.Flags().String("policy", "", "Set a IAM policy in JSON for the assumed credentials")
	envCmd.Flags().BoolP("mfa", "m", false, "Use MFA when assuming role")
	envCmd.Flags().String("mfacode", "", "Code to use for MFA")
	//revive:enable:line-length-limit
}
