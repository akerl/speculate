package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// AddAssumeFlags adds the flags for role assumption to a command object
func AddAssumeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("account", "a", "", "Account ID to assume role on (defaults to source account")
	cmd.Flags().StringP("session", "s", "", "Set session name for assumed role (defaults to origin user name)")
	AddMfaFlags(cmd)
}

// AssumeRole handles role assumption using
func AssumeRole(role, accountID, sessionName string) (Role, error) {
	arn, err := RoleArn(role, accountID)
	if err != nil {
		return Role{}, err
	}
	sessionName, err = SessionName(sessionName)
	if err != nil {
		return Role{}, err
	}

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(sessionName),
	}
	resp, err := StsSession.AssumeRole(params)
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
