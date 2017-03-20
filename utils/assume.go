package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AddAssumeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("account", "a", "", "Account ID to assume role on (defaults to source account")
	AddMfaFlags(cmd)
}

func AssumeRole(role string, flags *pflag.FlagSet) (Role, error) {
	arn, err := RoleArn(role)
	if err != nil {
		return Role{}, err
	}
	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String("roleSessionNameType"),
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
