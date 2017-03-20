package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

func AddAssumeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("account", "a", "", "Account ID to assume role on (defaults to source account")
	cmd.Flags().StringP("session", "s", "", "Set session name for assumed role (defaults to origin user name)")
	AddMfaFlags(cmd)
}

func AssumeRole(role, account_id, session_name string) (Role, error) {
	arn, err := RoleArn(role, account_id)
	if err != nil {
		return Role{}, err
	}
	session_name, err = SessionName(session_name)
	if err != nil {
		return Role{}, err
	}

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(session_name),
	}
	resp, err := StsSession.AssumeRole(params)
	if err != nil {
		return Role{}, err
	}

	creds := resp.Credentials
	new_role := Role{
		AccessKey:    *creds.AccessKeyId,
		SecretKey:    *creds.SecretAccessKey,
		SessionToken: *creds.SessionToken,
	}
	return new_role, nil
}
