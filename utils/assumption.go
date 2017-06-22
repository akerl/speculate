package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// Assumption describes the parameters that result in a Role
type Assumption struct {
	RoleName    string
	AccountID   string
	SessionName string
	Policy      string
	Lifetime
	Mfa
}

// ParseFlags for assumption object
func (a *Assumption) ParseFlags(cmd *cobra.Command) error {
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
	a.Policy, err = flags.GetString("policy")
	if err != nil {
		return err
	}
	err = a.parseLifetimeFlags(cmd)
	if err != nil {
		return err
	}
	err = a.parseMfaFlags(cmd)
	if err != nil {
		return err
	}
	return nil
}

// Execute actions a role assumption object with creds from the environment
func (a *Assumption) Execute() (Creds, error) {
	return a.ExecuteWithCreds(Creds{})
}

// ExecuteWithCreds actions a role assumption with provided creds
func (a *Assumption) ExecuteWithCreds(c Creds) (Creds, error) {
	newCreds := Creds{}
	arn, err := API.RoleArn(a.RoleName, a.AccountID)
	if err != nil {
		return newCreds, err
	}
	if a.SessionName == "" {
		a.SessionName, err = API.SessionName()
		if err != nil {
			return newCreds, err
		}
	}

	if a.LifetimeInt == 0 {
		a.LifetimeInt = 3600
	}

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(a.SessionName),
		DurationSeconds: aws.Int64(a.LifetimeInt),
	}
	if a.Policy != "" {
		params.Policy = aws.String(a.Policy)
	}
	if err := a.configureMfa(params); err != nil {
		return newCreds, err
	}

	client := API.ClientWithCreds(c)
	resp, err := client.AssumeRole(params)
	if err != nil {
		return newCreds, err
	}

	err = newCreds.NewFromStsSdk(resp.Credentials)
	return newCreds, err
}
