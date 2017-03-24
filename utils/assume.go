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
func AssumeRole(role string, accountID string, sessName string, useMfa bool, mfaCode string) (Role, error) {
	arn, err := roleArn(role, accountID)
	if err != nil {
		return Role{}, err
	}
	sessName, err = sessionName(sessName)
	if err != nil {
		return Role{}, err
	}

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(sessName),
	}
	if useMfa || mfaCode != "" {
		*params.SerialNumber, err = mfaArn()
		if err != nil {
			return Role{}, err
		}
		if mfaCode == "" {
			mfaCode, err = promptForMfa()
			if err != nil {
				return Role{}, err
			}
		}
		*params.TokenCode = mfaCode
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
