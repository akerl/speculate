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

// Execute actions a role assumption object
func (a *Assumption) Execute() (Creds, error) {
	creds := Creds{}
	arn, err := API.RoleArn(a.RoleName, a.AccountID)
	if err != nil {
		return creds, err
	}
	if a.SessionName == "" {
		a.SessionName, err = API.SessionName()
		if err != nil {
			return creds, err
		}
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
		return creds, err
	}

	client := API.Client()
	resp, err := client.AssumeRole(params)
	if err != nil {
		return creds, err
	}

	respCreds := resp.Credentials
	creds.New(map[string]string{
		"AccessKey":    *respCreds.AccessKeyId,
		"SecretKey":    *respCreds.SecretAccessKey,
		"SessionToken": *respCreds.SessionToken,
	})
	return creds, nil
}
