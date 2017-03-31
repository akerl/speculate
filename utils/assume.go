package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// Assumption describes the parameters that result in a Role
type Assumption struct {
	RoleName    string
	AccountID   string
	SessionName string
	UseMfa      bool
	MfaCode     string
	Lifetime    int64
	Policy      string
}

// AddAssumeFlags adds the flags for role assumption to a command object
func AddAssumeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("account", "a", "", "Account ID to assume role on (defaults to source account")
	cmd.Flags().StringP("session", "s", "", "Set session name for assumed role (defaults to origin user name)")
	cmd.Flags().Int64P("lifetime", "l", 3600, "Set lifetime of credentials in seconds (defaults to 3600 seconds / 1 hour, min 900, max 3600)")
	cmd.Flags().String("policy", "", "Set a IAM policy in JSON for the assumed credentials")
	cmd.Flags().BoolP("mfa", "m", false, "Use MFA when assuming role")
	cmd.Flags().String("mfacode", "", "Code to use for MFA")
}

// ParseFlags for assumption object
func (a Assumption) ParseFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()
	var err error
	a.AccountID, err = flags.GetString("account")
	if err != nil {
		return err
	}
	a.SessionName, err = flags.GetString("session")
	if err != nil {
		return err
	}
	a.Lifetime, err = flags.GetInt64("lifetime")
	if err != nil {
		return err
	}
	fmt.Printf("A Lifetime is %#v\n", a.Lifetime)
	if a.Lifetime < 900 || a.Lifetime > 3600 {
		return fmt.Errorf("Lifetime must be between 900 and 3600: %d", a.Lifetime)
	}
	fmt.Printf("B Lifetime is %#v\n", a.Lifetime)
	a.Policy, err = flags.GetString("policy")
	if err != nil {
		return err
	}
	a.UseMfa, err = flags.GetBool("mfa")
	if err != nil {
		return err
	}
	a.MfaCode, err = flags.GetString("mfacode")
	if err != nil {
		return err
	}
	return nil
}

// AssumeRole actions a role assumption object
func (a Assumption) AssumeRole() (Role, error) {
	arn, err := API.RoleArn(a.RoleName, a.AccountID)
	if err != nil {
		return Role{}, err
	}
	if a.SessionName == "" {
		a.SessionName, err = API.SessionName()
		if err != nil {
			return Role{}, err
		}
	}

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(a.SessionName),
		DurationSeconds: aws.Int64(a.Lifetime),
	}
	if a.Policy != "" {
		params.Policy = aws.String(a.Policy)
	}
	if err := configureMfa(a, params); err != nil {
		return Role{}, err
	}

	client := API.Client()
	resp, err := client.AssumeRole(params)
	if err != nil {
		return Role{}, err
	}

	creds := resp.Credentials
	newRole := Role{
		AccessKey:    *creds.AccessKeyId,
		SecretKey:    *creds.SecretAccessKey,
		SessionToken: *creds.SessionToken,
	}
	return newRole, nil
}
